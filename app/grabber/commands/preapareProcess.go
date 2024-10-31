package commands

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/grabber/commands/grabberflag"
	"github.com/Galdoba/ffstuff/app/grabber/internal/actions"
	"github.com/Galdoba/ffstuff/app/grabber/internal/copyprocess"
	"github.com/Galdoba/ffstuff/app/grabber/internal/origin"
	"github.com/Galdoba/ffstuff/app/grabber/internal/process"
	"github.com/Galdoba/ffstuff/app/grabber/internal/sourcesort"
	"github.com/Galdoba/ffstuff/app/grabber/internal/target"
	logman "github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

func preapareProcess(c *cli.Context) (copyprocess.CopyProcess, error) {
	logman.Debug(logman.NewMessage("start process preparation"))
	//Setup process
	logman.Debug(logman.NewMessage("check flags"))
	if err := grabberflag.ValidateGrabFlags(c); err != nil {
		return nil, logman.Errorf("flag validation failed: %v", err)
	}
	logman.Debug(logman.NewMessage("check arguments"))
	logman.Debug(logman.NewMessage("set process options"))
	options := process.DefineGrabOptions(c, cfg)
	process, err := process.New(options...)
	if err != nil {
		return nil, logman.Errorf("process creation failed")
	}
	//Setup sources
	if err := origin.ConstructorSetup(
		origin.WithFilePriority(cfg.FILE_PRIORITY_WEIGHTS),
		origin.WithDirectoryPriority(cfg.DIRECTORY_PRIORITY_WEIGHTS),
		origin.KillSignal(process.DeleteDecidion),
		origin.WithMarkerExt(cfg.MARKER_FILE_EXTENTION),
	); err != nil {
		return nil, logman.Errorf("source constructor setup failed: %v", err)
	}
	//Setup target manager
	tm, err := target.NewTargetManager(cfg,
		target.WithDestination(process.DestinationDir),
		target.WithCopyHandling(process.CopyDecidion))
	if err != nil {
		return nil, logman.Errorf("target manager setup failed: %v", err)
	}

	//compile source list
	sources := []origin.Origin{}
	for grNum, arg := range c.Args().Slice() {
		gr := fmt.Sprintf("group_%02d", grNum)
		src := origin.New(arg, gr)
		sources = append(sources, src)
		related, err := actions.DiscoverRelatedFiles(src)
		if err != nil {
			logman.Warn("failed to discover related files: %v", err)
		}
		for _, found := range related {
			sources = append(sources, origin.New(found, gr))
		}
	}
	logman.Printf("%v source files received", len(sources))

	//Sort
	sources, err = sourcesort.Sort(process, sources...)
	if err != nil {
		return nil, logman.Errorf("sort error: %v", err)
	}
	//targeting
	filteredSources := []origin.Origin{}
	for _, src := range sources {
		tgtName, err := tm.NewTarget(src)
		if err != nil {
			logman.Errorf("failed to create target for source '%v': %v", src.Name())
			continue
		}
		if tgtName == "" {
			switch process.CopyDecidion {
			case grabberflag.VALUE_COPY_SKIP:
				logman.Warn("skip %v", src.Name())
				continue
			default:
				logman.Errorf("failed to compile target for source '%v'", src.Name())
				continue
			}
		}
		process.SourceTargetMap[src] = tgtName
		filteredSources = append(filteredSources, src)
	}
	copyProc := copyprocess.NewCopyAction(process.SourceTargetMap,
		copyprocess.WithDestination(process.DestinationDir),
		copyprocess.WithMarkerExt(cfg.MARKER_FILE_EXTENTION),
		copyprocess.WithDeleteDecidion(process.DeleteDecidion),
		copyprocess.WithSourcePaths(filteredSources...),
	)
	logman.Debug(logman.NewMessage("end process preparation"))
	return copyProc, nil
}
