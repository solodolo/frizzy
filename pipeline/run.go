package pipeline

import (
	"context"
	"log"
	"os"
	"sync"
)

func RunPipeline(pathChan <-chan string, handler func(context.Context, *os.File) []<-chan error) error {
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

		errChans = append(errChans, handler(ctx, f)...)
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

	merge := func(idx int, ec <-chan error, errChans []<-chan error) {
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
		go merge(i, ec, errChans)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	return errChan
}
