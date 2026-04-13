import { registerAs } from '@nestjs/config';

export default registerAs('database', () => ({
  type: (process.env.DB_CONNECTION ?? 'sqlite') as 'mysql' | 'sqlite',
  host: process.env.DB_HOST ?? '127.0.0.1',
  port: Number.parseInt(process.env.DB_PORT ?? '3306', 10),
  database: process.env.DB_DATABASE ?? 'incharge',
  username: process.env.DB_USERNAME ?? 'root',
  password: process.env.DB_PASSWORD ?? '',
}));
