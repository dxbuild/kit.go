package rt

import (
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

func ProvideLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return logger
}

type EventLogger struct {
	log zerolog.Logger
}

func (el *EventLogger) LogEvent(event fxevent.Event) {
	log := el.log.With().Str("cat", "FX").Logger()
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		log.Info().Msgf("HOOK OnStart\t\t%s executing (caller: %s)", e.FunctionName, e.CallerName)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			log.Error().Err(e.Err).Msgf("HOOK OnStart\t\t%s called by %s failed in %s", e.FunctionName, e.CallerName, e.Runtime)
		} else {
			log.Info().Msgf("HOOK OnStart\t\t%s called by %s ran successfully in %s", e.FunctionName, e.CallerName, e.Runtime)
		}
	case *fxevent.OnStopExecuting:
		log.Info().Msgf("HOOK OnStop\t\t%s executing (caller: %s)", e.FunctionName, e.CallerName)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			log.Error().Err(e.Err).Msgf("HOOK OnStop\t\t%s called by %s failed in %s", e.FunctionName, e.CallerName, e.Runtime)
		} else {
			log.Info().Msgf("HOOK OnStop\t\t%s called by %s ran successfully in %s", e.FunctionName, e.CallerName, e.Runtime)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			log.Error().Err(e.Err).Msgf("ERROR\tFailed to supply %v", e.TypeName)
		}
	case *fxevent.Provided:
		var privateStr string
		if e.Private {
			privateStr = " (PRIVATE)"
		}
		for _, rtype := range e.OutputTypeNames {
			log.Info().Msgf("PROVIDE%v\t%v <= %v", privateStr, rtype, e.ConstructorName)
		}
		if e.Err != nil {
			log.Error().Err(e.Err).Msgf("Error after options were applied")
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			log.Info().Msgf("REPLACE\t%v", rtype)
		}
		if e.Err != nil {
			log.Error().Err(e.Err).Msgf("ERROR\tFailed to replace")
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				log.Info().Msgf("DECORATE\t%v <= %v from module %q", rtype, e.DecoratorName, e.ModuleName)
			} else {
				log.Info().Msgf("DECORATE\t%v <= %v", rtype, e.DecoratorName)
			}
		}
		if e.Err != nil {
			log.Error().Err(e.Err).Msgf("ERROR\tFailed to decorate")
		}
	case *fxevent.Run:
		var moduleStr string
		if e.ModuleName != "" {
			moduleStr = " from module " + e.ModuleName
		}
		log.Info().Msgf("RUN\t%v: %v%v", e.Kind, e.Name, moduleStr)
		if e.Err != nil {
			log.Error().Err(e.Err).Msgf("ERROR\tFailed to run")
		}
	case *fxevent.Invoking:
		if e.ModuleName != "" {
			log.Info().Msgf("INVOKE\t\t%s from module %q", e.FunctionName, e.ModuleName)
		} else {
			log.Info().Msgf("INVOKE\t\t%s", e.FunctionName)
		}
	case *fxevent.Invoked:
		if e.Err != nil {
			log.Error().Err(e.Err).Msgf("ERROR\tFailed to invoke")
		}
	case *fxevent.Stopping:
		log.Info().Msgf("STOPPING\t\t%s", e.Signal.String())
	case *fxevent.Stopped:
		if e.Err != nil {
			log.Error().Err(e.Err).Msgf("ERROR\tFailed to stop")
		}
	case *fxevent.RollingBack:
		log.Info().Msgf("ROLLING BACK\t\t%s", e.StartErr)
	case *fxevent.RolledBack:
		if e.Err != nil {
			log.Error().Err(e.Err).Msgf("ERROR\tFailed to rollback")
		}
	case *fxevent.Started:
		if e.Err != nil {
			log.Error().Err(e.Err).Msgf("ERROR\tFailed to start")
		} else {
			log.Info().Msgf("STARTED")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			log.Error().Err(e.Err).Msgf("ERROR\tFailed to initialize logger")
		} else {
			log.Info().Msgf("LOGGER INITIALIZED")
		}
	}
}

func WithLogger(log zerolog.Logger) fxevent.Logger {
	return &EventLogger{log: log}
}
