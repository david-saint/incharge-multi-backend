# InCharge Server

Production-grade backend implementation of the `SPEC.md` contract in `Node.js`, `TypeScript`, `NestJS`, and `MySQL`/`SQLite`.

## Stack

- Node.js 22+
- NestJS 11
- TypeORM
- MySQL 8 for production
- SQLite in-memory for local smoke tests and e2e
- JWT auth for API users
- Session-cookie auth for admins
- Nodemailer for verification and password reset emails
- Vue 3 browser shell for admin entry pages

## Implemented Surface

### API

- `GET /api/v1/global`
- `GET /api/v1/global/contraception-reasons`
- `GET /api/v1/global/contraception-reasons/:id`
- `GET /api/v1/global/education-levels`
- `GET /api/v1/global/faq-groups`
- `GET /api/v1/global/faq-groups/:id`
- `POST /api/v1/user/register`
- `POST /api/v1/user/login`
- `POST /api/v1/user/logout`
- `GET /api/v1/user/refresh`
- `POST /api/v1/user/password/email`
- `POST /api/v1/user/password/reset`
- `GET /api/v1/user/email/verify/:id`
- `GET /api/v1/user/email/resend`
- `GET /api/v1/user/email/success`
- `GET /api/v1/user`
- `POST /api/v1/user/profile`
- `GET /api/v1/user/profile`
- `POST /api/v1/user/profile/algorithm`
- `GET /api/v1/user/clinics`
- Legacy session-protected clinic/user management routes under `/api/v1/user/...`

### Admin Web

- `GET /`
- `GET /admin`
- `GET /loginView`
- `POST /login`
- `GET /logout`
- `GET /privacy`
- `GET /panel`
- `POST /admin`
- `GET /allAdmins`
- `GET /getAdminDet`
- `PUT /updateAdmin/:admin_id`
- `GET /getUsers`
- `DELETE /deleteUser/:user_id`
- `PUT /revertDeletedUser/:user_id`
- `GET /getDeletedUsers`
- `GET /getClinics`
- `POST /addClinic`
- `PUT /updateClinic/:clinic_id`
- `DELETE /deleteClinic/:clinic_id`
- `GET /getDeletedClinics`
- `PUT /revertDeletedClinic/:clinic_id`
- `GET /getContraceptionReason`
- `GET /getEducationalLevels`
- `GET /algo`
- `POST /algo`
- `PUT /algo/:id`

## Setup

1. Install dependencies.

```bash
npm install
```

2. Copy env defaults.

```bash
cp .env.example .env
```

3. Choose a database mode.

For local smoke testing:

```env
DB_CONNECTION=sqlite
DB_DATABASE=:memory:
DB_SYNCHRONIZE=true
AUTO_SEED=true
```

For MySQL:

```env
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=incharge
DB_USERNAME=root
DB_PASSWORD=secret
DB_SYNCHRONIZE=true
AUTO_SEED=true
SESSION_STORE=mysql
```

4. Start the server.

```bash
npm run start:dev
```

## Important Env Vars

- `APP_URL`: server base URL
- `APP_API_DOMAIN`: API domain used in generated verify links
- `APP_USER_DOMAIN`: frontend domain used in verification and reset links
- `JWT_SECRET`: JWT signing secret
- `JWT_TTL_SECONDS`: access-token TTL
- `SESSION_SECRET`: admin session secret
- `SESSION_STORE`: `memory` or `mysql`
- `MAIL_ENABLED`: `true` to send SMTP mail
- `PASSWORD_RESET_EXPIRY_MINUTES`: default `60`
- `EMAIL_VERIFICATION_EXPIRY_MINUTES`: default `4320`

## Reference Data

The app seeds the following on startup when `AUTO_SEED=true`:

- contraception reasons
- education levels
- FAQ groups and placeholder FAQ content
- countries, states, locations from CSVs in `database/data`
- algorithm rows from `database/sql/algorithms.sql`

Replace the sample CSV and SQL files with production data as needed.

## Admin UI

The admin entry pages are server-served HTML with a lightweight Vue 3 browser shell:

- `public/admin-login.html`
- `public/register-super-admin.html`
- `public/admin-panel.html`

They intentionally preserve the spec’s model of server-rendered entry pages backed by a Vue SPA shell.

## Security and Behavior Notes

- Open CORS with `Authorization` exposed
- Global throttling: `120 req/min`
- Email verification and resend throttling: `6 req/min`
- JWT responses include both JSON token body and `Authorization` response header
- User logout revokes the active JWT session
- Admin web routes use cookie-backed sessions
- Verified-admin-only protection is enforced on panel and protected web routes
- Soft deletes are implemented on all spec-required entities

## Testing

Build:

```bash
npm run build
```

E2E:

```bash
npm run test:e2e
```

Current e2e coverage verifies:

- global API health/reference endpoints
- user register/login/profile/algorithm flow
- admin bootstrap/login/clinic management flow
- email verification signature handling
- password reset email request flow
- verified-admin route protection

## Gaps and Follow-up Work

This implementation is feature-complete for the tested contract, but a production rollout should still add:

1. Real Vue build pipeline instead of CDN-loaded browser Vue.
2. Proper TypeORM migrations instead of `synchronize` for production deployments.
3. Full production seed files for algorithms and geography.
4. Broader e2e coverage for every legacy route variation and edge-case validation message.
5. CSRF protection if the admin panel evolves beyond same-origin session posting.
