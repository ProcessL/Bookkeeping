package core

import (
	"github.com/dotdancer/gogofly/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

func InitLogger() *zap.SugaredLogger {

	logMode := zapcore.DebugLevel
	if !global.Config.Zap.Develop {
		logMode = zapcore.InfoLevel
	}

	core := zapcore.NewCore(getEncoder(), zapcore.NewMultiWriteSyncer(getWriteSyncer(), zapcore.AddSync(os.Stdout)), logMode)
	return zap.New(core, zap.AddCaller()).Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Local().Format(time.DateTime))
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getWriteSyncer() zapcore.WriteSyncer {
	stSeparator := string(filepath.Separator)
	stRootDir, _ := os.Getwd()
	stLogFilePath := stRootDir + stSeparator + "log" + stSeparator + time.Now().Format(time.DateOnly) + ".txt"
	lumberJackSyncer := &lumberjack.Logger{
		Filename:   stLogFilePath,
		MaxSize:    global.Config.Zap.MaxSize,
		MaxBackups: global.Config.Zap.MaxBackups,
		MaxAge:     global.Config.Zap.MaxAge,
		Compress:   false, // disabled by default
	}
	return zapcore.AddSync(lumberJackSyncer)
}
