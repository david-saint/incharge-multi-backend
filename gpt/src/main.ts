import session from 'express-session';
import { UnprocessableEntityException, ValidationPipe } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { NestFactory } from '@nestjs/core';
import cookieParser from 'cookie-parser';
import helmet from 'helmet';
import type { Store } from 'express-session';
import { AppModule } from './app.module';
import { HttpExceptionFilter } from './common/filters/http-exception.filter';
import { formatValidationErrors } from './common/utils/validation.util';

async function bootstrap() {
  const app = await NestFactory.create(AppModule, { bufferLogs: true });
  const config = app.get(ConfigService);

  app.setGlobalPrefix('');
  app.enableCors({
    origin: '*',
    methods: '*',
    allowedHeaders: '*',
    exposedHeaders: ['Authorization'],
  });
  app.use(helmet({ crossOriginResourcePolicy: false }));
  app.use(cookieParser());
  app.use(
    session({
      name: 'incharge_admin',
      secret: config.get<string>('app.sessionSecret', 'incharge-session-secret'),
      resave: false,
      saveUninitialized: false,
      rolling: true,
      store: (await createSessionStore(config)) as Store | undefined,
      cookie: {
        httpOnly: true,
        sameSite: 'lax',
        secure: config.get<string>('app.env') === 'production',
        maxAge: 1000 * 60 * 60 * 24 * 7,
      },
    }),
  );

  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      transform: true,
      exceptionFactory: (errors) =>
        new UnprocessableEntityException({
          errors: formatValidationErrors(errors),
        }),
    }),
  );
  app.useGlobalFilters(new HttpExceptionFilter(config));
  app.getHttpAdapter().getInstance().set('trust proxy', 1);

  await app.listen(config.get<number>('app.port', 3000));
}

async function createSessionStore(
  config: ConfigService,
): Promise<Store | undefined> {
  const sessionStore = config.get<string>('app.sessionStore', 'memory');
  const databaseType = config.get<string>('database.type', 'sqlite');
  if (sessionStore !== 'mysql' || databaseType !== 'mysql') {
    return undefined;
  }

  const mysqlSession = await import('express-mysql-session');
  const sessionStoreFactory = mysqlSession.default(session);
  return new sessionStoreFactory({
    host: config.get<string>('database.host', '127.0.0.1'),
    port: config.get<number>('database.port', 3306),
    user: config.get<string>('database.username', 'root'),
    password: config.get<string>('database.password', ''),
    database: config.get<string>('database.database', 'incharge'),
    clearExpired: true,
    createDatabaseTable: true,
    schema: {
      tableName: 'admin_sessions',
    },
  });
}

void bootstrap();
