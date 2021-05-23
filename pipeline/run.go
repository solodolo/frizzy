package pipeline

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func WalkFiles(inputPath string) (<-chan string, <-chan error) {
	pathChan := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(pathChan)
		defer close(errChan)

		walkErr := filepath.WalkDir(inputPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			// ignore nested dirs for now
			if d.IsDir() {
				return nil
			}

			pathChan <- path
			return nil
		})

		if walkErr != nil {
			errChan <- walkErr
		}
	}()

	return pathChan, errChan
}

func RunPipeline(pathChan <-chan string, handler func(context.Context, *os.File) <-chan error) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChans := []<-chan error{}
	files := make([]*os.File, len(pathChan))

	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	for inputPath := range pathChan {
		log.Printf("    %s", inputPath)

		f, err := os.Open(inputPath)
		files = append(files, f)

		if err != nil {
			log.Printf("    pipeline error: failed to open %s, %s\n", inputPath, err)
			continue
		}

		errChans = append(errChans, handler(ctx, f))
	}

	errChan := mergeErrChans(ctx, errChans)

	for err := range errChan {
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func mergeErrChans(ctx context.Context, errChans []<-chan error) <-chan error {
	wg := sync.WaitGroup{}
	wg.Add(len(errChans))

	errChan := make(chan error)

	merge := func(idx int, ec <-chan error) {
		defer wg.Done()
		for err := range ec {
			select {
			case errChan <- err:
			case <-ctx.Done():
				return
			}
		}
	}

	for i, ec := range errChans {
		go merge(i, ec)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	return errChan
}
