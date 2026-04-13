import { registerAs } from '@nestjs/config';

export default registerAs('app', () => ({
  env: process.env.APP_ENV ?? 'local',
  port: Number.parseInt(process.env.PORT ?? '3000', 10),
  url: process.env.APP_URL ?? 'http://localhost:3000',
  apiDomain:
    process.env.APP_API_DOMAIN ?? process.env.APP_URL ?? 'http://localhost:3000',
  userDomain: process.env.APP_USER_DOMAIN ?? 'http://localhost:5173',
  jwtSecret: process.env.JWT_SECRET ?? 'change-me',
  jwtTtlSeconds: Number.parseInt(process.env.JWT_TTL_SECONDS ?? '3600', 10),
  dbSynchronize: ['1', 'true', 'yes'].includes(
    (process.env.DB_SYNCHRONIZE ?? 'true').toLowerCase(),
  ),
  autoSeed: ['1', 'true', 'yes'].includes(
    (process.env.AUTO_SEED ??
      (process.env.NODE_ENV === 'test' ? 'true' : 'false')).toLowerCase(),
  ),
  sessionSecret: process.env.SESSION_SECRET ?? 'incharge-session-secret',
  sessionStore: process.env.SESSION_STORE ?? 'memory',
  passwordResetExpiryMinutes: Number.parseInt(
    process.env.PASSWORD_RESET_EXPIRY_MINUTES ?? '60',
    10,
  ),
  emailVerificationExpiryMinutes: Number.parseInt(
    process.env.EMAIL_VERIFICATION_EXPIRY_MINUTES ?? '4320',
    10,
  ),
  mailEnabled: ['1', 'true', 'yes'].includes(
    (process.env.MAIL_ENABLED ?? 'false').toLowerCase(),
  ),
}));
