import 'express-session';

declare module 'express-session' {
  interface SessionData {
    adminId?: number;
  }
}

declare module 'express-serve-static-core' {
  interface Request {
    user?: unknown;
  }
}
