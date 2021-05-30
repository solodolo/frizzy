package main

import (
	"flag"
	"log"
	"os"

	"mettlach.codes/frizzy/config"
	"mettlach.codes/frizzy/file"
	"mettlach.codes/frizzy/pipeline"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	startDevServer := flag.Bool("d", false, "start a web server to serve files in output directory")
	devServerPort := flag.Int("p", 8080, "the port the web server will listen on")
	clearOutput := flag.Bool("c", false, "clear any existing output")
	flag.Parse()

	configPath := os.Args[len(os.Args)-1]

	if config, err := config.LoadConfig(configPath); err != nil {
		log.Printf("error loading config: %s\n", err)
		return
	} else {
		if *clearOutput {
			log.Printf("removing %s\n", config.OutputPath)
			if err := clearOutputDirectory(config.OutputPath); err != nil {
				log.Print(err)
				return
			}
		}

		templatePathChan, _ := pipeline.WalkFiles(config.GetTemplatePath())
		contentPathChan, _ := pipeline.WalkFiles(config.GetContentPath())
		pagesPathChan, _ := pipeline.WalkFiles(config.GetPagesPath())
		// have to process templates first, then content, then pages
		log.Println("pipelining template files")
		if err := pipeline.RunPipeline(templatePathChan, pipeline.TemplateCacheHandler); err != nil {
			log.Println("exiting")
			return
		} else {
			log.Println("finished template files")
		}

		log.Println("pipelining content files")
		if err := pipeline.RunPipeline(contentPathChan, pipeline.FullPipelineNullRenderer); err != nil {
			log.Println("exiting")
			return
		} else {
			log.Println("finished content files")
		}

		log.Println("pipelining page files")
		if err := pipeline.RunPipeline(pagesPathChan, pipeline.FullPipelineNullRenderer); err != nil {
			log.Println("exiting")
			return
		} else {
			log.Println("finished page files")
		}

		if *startDevServer {
			log.Println("starting development server...")
			server := file.DevServer{ServerRoot: config.OutputPath, Port: *devServerPort}
			server.ListenAndServe()
		}

		log.Println("Done")
	}
}

func printUsage() {
	log.Println("usage: frizzy [-c] /path/to/config.json")
}

func clearOutputDirectory(outputDir string) error {
	return os.RemoveAll(outputDir)
}
