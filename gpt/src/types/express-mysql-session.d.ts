declare module 'express-mysql-session' {
  import type session from 'express-session';

  interface Options {
    host?: string;
    port?: number;
    user?: string;
    password?: string;
    database?: string;
    clearExpired?: boolean;
    createDatabaseTable?: boolean;
    schema?: {
      tableName?: string;
    };
  }

  export default function expressMysqlSession(
    sessionModule: typeof session,
  ): new (options: Options) => session.Store;
}
