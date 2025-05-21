import winston, { format } from 'winston';

const logger = winston.createLogger({
    level: 'info',
    format: winston.format.combine(
        winston.format.json(),
        winston.format.colorize({ level: true }),
        winston.format.timestamp({ format: 'YYYY-MM-DD HH:mm:ss,SSS' }),
        winston.format.printf(({ level, timestamp, message }) => `${timestamp} | ${level} | ${message}`),
    ),
    transports: [
        new winston.transports.Console(),
    ],
});

export default logger;