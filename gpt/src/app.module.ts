import { join } from 'node:path';
import { Module } from '@nestjs/common';
import { APP_GUARD } from '@nestjs/core';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { ServeStaticModule } from '@nestjs/serve-static';
import { ThrottlerGuard, ThrottlerModule } from '@nestjs/throttler';
import { TypeOrmModule } from '@nestjs/typeorm';
import { LoggerModule } from 'nestjs-pino';
import { AdminSessionGuard } from './common/guards/admin-session.guard';
import { VerifiedAdminGuard } from './common/guards/verified-admin.guard';
import appConfig from './config/app.config';
import databaseConfig from './config/database.config';
import { envValidationSchema } from './config/env.validation';
import { AdminModule } from './admin/admin.module';
import { AlgorithmsModule } from './algorithms/algorithms.module';
import { AuthModule } from './auth/auth.module';
import { ClinicsModule } from './clinics/clinics.module';
import { DatabaseModule } from './database/database.module';
import { GlobalApiModule } from './global/global.module';
import { MailModule } from './mail/mail.module';
import { ProfilesModule } from './profiles/profiles.module';
import { UsersModule } from './users/users.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      load: [appConfig, databaseConfig],
      validationSchema: envValidationSchema,
    }),
    LoggerModule.forRoot({
      pinoHttp: {
        level: process.env.APP_ENV === 'production' ? 'info' : 'debug',
        transport:
          process.env.APP_ENV === 'production'
            ? undefined
            : { target: 'pino-pretty', options: { singleLine: true } },
      },
    }),
    ServeStaticModule.forRoot({
      rootPath: join(__dirname, '..', 'public'),
      serveRoot: '/static',
    }),
    ThrottlerModule.forRoot([
      {
        ttl: 60_000,
        limit: 120,
      },
    ]),
    TypeOrmModule.forRootAsync({
      inject: [ConfigService],
      useFactory: (config: ConfigService) => {
        const type = config.get<'mysql' | 'sqlite'>('database.type', 'sqlite');
        return {
          type,
          database:
            type === 'sqlite'
              ? config.get<string>('database.database', ':memory:')
              : config.get<string>('database.database', 'incharge'),
          host: type === 'mysql' ? config.get<string>('database.host') : undefined,
          port:
            type === 'mysql' ? config.get<number>('database.port', 3306) : undefined,
          username:
            type === 'mysql' ? config.get<string>('database.username') : undefined,
          password:
            type === 'mysql' ? config.get<string>('database.password') : undefined,
          synchronize: config.get<boolean>('app.dbSynchronize', false),
          autoLoadEntities: true,
          logging: config.get<string>('app.env') === 'local',
        };
      },
    }),
    DatabaseModule,
    MailModule,
    AuthModule,
    GlobalApiModule,
    ProfilesModule,
    ClinicsModule,
    AlgorithmsModule,
    UsersModule,
    AdminModule,
  ],
  providers: [
    AdminSessionGuard,
    VerifiedAdminGuard,
    {
      provide: APP_GUARD,
      useClass: ThrottlerGuard,
    },
  ],
})
export class AppModule {}
