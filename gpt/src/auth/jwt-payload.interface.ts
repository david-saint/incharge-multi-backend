export interface JwtPayload {
  sub: number;
  email: string;
  type: 'user';
  jti: string;
}
