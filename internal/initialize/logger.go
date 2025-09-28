package initialize

import (
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

func InitLogger() {
	loggerPAth := global.WORK_DIR + "/logs"
	if ok, _ := utils.PathExists(loggerPAth); !ok { // 判断是否有Director文件夹
		_ = os.Mkdir(loggerPAth, os.ModePerm)
	}
	core := getEncoderCore(loggerPAth + "/server.log")
	global.LOG = zap.New(core, zap.AddCaller())
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "test",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	return config
}

func getEncoder() zapcore.Encoder {
	zap.NewProductionEncoderConfig()
	encoderConfig := getEncoderConfig()
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getEncoderCore(fileName string) (core zapcore.Core) {
	writer := getWriteSyncer(fileName) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(getEncoder(), writer, getLevel())
}

func getLevel() zapcore.LevelEnabler {
	switch strings.ToUpper(global.CONFIG.Logger.Level) {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func getWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,  //日志文件的位置
		MaxSize:    10,    //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 200,   //保留旧文件的最大个数
		MaxAge:     30,    //保留旧文件的最大天数
		Compress:   false, //是否压缩/归档旧文件
	}
	if global.CONFIG.Logger.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05.000"))
}
