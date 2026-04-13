import session from 'express-session';
import { INestApplication, UnprocessableEntityException, ValidationPipe } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Test, TestingModule } from '@nestjs/testing';
import cookieParser from 'cookie-parser';
import request from 'supertest';
import { App } from 'supertest/types';
import { AppModule } from './../src/app.module';
import { HttpExceptionFilter } from '../src/common/filters/http-exception.filter';
import { signPayload } from '../src/common/utils/security.util';
import { formatValidationErrors } from '../src/common/utils/validation.util';

describe('InCharge Server (e2e)', () => {
  let app: INestApplication<App>;

  beforeAll(async () => {
    process.env.APP_ENV = 'test';
    process.env.DB_CONNECTION = 'sqlite';
    process.env.DB_DATABASE = ':memory:';
    process.env.DB_SYNCHRONIZE = 'true';
    process.env.AUTO_SEED = 'true';
    process.env.JWT_SECRET = 'test-secret-value';
    process.env.SESSION_SECRET = 'test-session-value';
    process.env.MAIL_ENABLED = 'false';

    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = moduleFixture.createNestApplication();
    const config = app.get(ConfigService);

    app.enableCors({
      origin: '*',
      methods: '*',
      allowedHeaders: '*',
      exposedHeaders: ['Authorization'],
    });
    app.use(cookieParser());
    app.use(
      session({
        name: 'incharge_admin',
        secret: config.get<string>('app.sessionSecret', 'test-session-value'),
        resave: false,
        saveUninitialized: false,
      }),
    );
    app.useGlobalPipes(
      new ValidationPipe({
        whitelist: true,
        transform: true,
        exceptionFactory: (errors) =>
          new UnprocessableEntityException({ errors: formatValidationErrors(errors) }),
      }),
    );
    app.useGlobalFilters(new HttpExceptionFilter(config));

    await app.init();
  });

  afterAll(async () => {
    await app.close();
  });

  it('returns the health check', async () => {
    await request(app.getHttpServer())
      .get('/api/v1/global')
      .expect(200)
      .expect('Hello, World!');
  });

  it('lists seeded reference data', async () => {
    const reasonsResponse = await request(app.getHttpServer())
      .get('/api/v1/global/contraception-reasons')
      .expect(200);

    expect(reasonsResponse.body).toHaveLength(3);

    const faqResponse = await request(app.getHttpServer())
      .get('/api/v1/global/faq-groups/1')
      .expect(200);

    expect(faqResponse.body.status).toBe('faq.get_content');
  });

  it('applies the default API limit and the stricter email verification limit', async () => {
    const globalResponse = await request(app.getHttpServer()).get('/api/v1/global').expect(200);

    expect(globalResponse.headers['x-ratelimit-limit']).toBe('120');

    await request(app.getHttpServer())
      .post('/api/v1/user/register')
      .send({
        name: 'Rate Limit User',
        email: 'ratelimit@example.com',
        password: 'secret123',
      })
      .expect(201);

    const config = app.get(ConfigService);
    const expires = String(Date.now() + 60_000);
    const signaturePayload = `1:${expires}`;
    const signature = signPayload(signaturePayload, config.get<string>('app.jwtSecret', ''));

    const verifyResponse = await request(app.getHttpServer())
      .get(`/api/v1/user/email/verify/1?expires=${expires}&signature=${signature}`)
      .expect(302);

    expect(verifyResponse.headers['x-ratelimit-limit']).toBe('6');
  });

  it('registers, authenticates, creates a profile, and stores an algorithm plan', async () => {
    const registerResponse = await request(app.getHttpServer())
      .post('/api/v1/user/register')
      .send({
        name: 'Jane Doe',
        email: 'jane@example.com',
        phone: '+12345678901',
        password: 'secret123',
      })
      .expect(201);

    expect(registerResponse.body.status).toBe(true);

    const loginResponse = await request(app.getHttpServer())
      .post('/api/v1/user/login')
      .send({ email: 'jane@example.com', password: 'secret123' })
      .expect(200);

    const token = loginResponse.body.token as string;
    expect(loginResponse.headers.authorization).toMatch(/^Bearer /);

    await request(app.getHttpServer())
      .post('/api/v1/user/profile')
      .set('Authorization', `Bearer ${token}`)
      .send({
        gender: 'FEMALE',
        age: 25,
        address: '12 Main Street',
        sexually_active: true,
      })
      .expect(201);

    const profileResponse = await request(app.getHttpServer())
      .get('/api/v1/user/profile?with=reason,educationLevel')
      .set('Authorization', `Bearer ${token}`)
      .expect(200);

    expect(profileResponse.body.gender).toBe('FEMALE');
    expect(profileResponse.body.reason.value).toBeDefined();

    await request(app.getHttpServer())
      .post('/api/v1/user/profile/algorithm')
      .set('Authorization', `Bearer ${token}`)
      .send({ plan: 'Progestin Only Pills' })
      .expect(200);

    const currentUserResponse = await request(app.getHttpServer())
      .get('/api/v1/user')
      .set('Authorization', `Bearer ${token}`)
      .expect(200);

    expect(currentUserResponse.body.id).toBeGreaterThan(0);
  });

  it('creates the first super admin, logs in, and manages clinics', async () => {
    const agent = request.agent(app.getHttpServer());

    const createAdminResponse = await agent.post('/admin').send({
      firstname: 'Super',
      lastname: 'Admin',
      phone: '+12345678902',
      email: 'admin@example.com',
      verified: 'Y',
      userType: 'Super',
      password: 'secret123',
    });

    expect(createAdminResponse.status).toBe(201);

    await agent
      .post('/login')
      .send({ email: 'admin@example.com', password: 'secret123' })
      .expect(200);

    const clinicResponse = await agent.post('/addClinic').send({
      name: 'Lagos Family Clinic',
      address: '1 Clinic Way',
      latitude: 6.5244,
      longitude: 3.3792,
      added_by_id: 1,
    });

    expect(clinicResponse.status).toBe(201);

    const publicClinicResponse = await request(app.getHttpServer())
      .get('/api/v1/user/clinics?latitude=6.5244&longitude=3.3792&radius=50&mode=km')
      .expect(200);

    expect(Array.isArray(publicClinicResponse.body)).toBe(true);
    expect(publicClinicResponse.body[0].distance).toContain('km');

    await agent.delete('/deleteClinic/1').expect(200);
    await agent.put('/revertDeletedClinic/1').expect(200);
  });

  it('verifies email signatures and password reset flow', async () => {
    await request(app.getHttpServer())
      .post('/api/v1/user/register')
      .send({
        name: 'John Doe',
        email: 'john@example.com',
        password: 'secret123',
      })
      .expect(201);

    const config = app.get(ConfigService);
    const expires = String(Date.now() + 60_000);
    const signaturePayload = `2:${expires}`;
    const signature = signPayload(signaturePayload, config.get<string>('app.jwtSecret', ''));

    await request(app.getHttpServer())
      .get(`/api/v1/user/email/verify/2?expires=${expires}&signature=${signature}`)
      .expect(302);

    await request(app.getHttpServer())
      .post('/api/v1/user/password/email')
      .send({ email: 'john@example.com' })
      .expect(200);
  });

  it('protects verified-admin-only web routes', async () => {
    await request(app.getHttpServer()).get('/panel').expect(401);
    await request(app.getHttpServer()).get('/getUsers').expect(401);
  });
});
