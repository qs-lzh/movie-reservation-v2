Use Zap to log instead of gin's logger

There're three forms of log:
- ZapLogger middleware logs every http requests at INFO level automatically
- ErrorLogger middleware logs detailed error of the failed http requests
- The zap logger also be wrapped in App struct, so handlers can use it to log business-level information

ErrorLogger examines the errors in context written by handler, and log them. Handler uses ctx.Error(err) to add the error into context. This reduces the repeated error-logging code in handlers.
