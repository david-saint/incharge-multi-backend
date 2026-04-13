import {
  Body,
  Controller,
  Get,
  HttpException,
  HttpCode,
  Post,
  Param,
  ParseIntPipe,
  Query,
  Redirect,
  Req,
  Res,
  UseGuards,
} from '@nestjs/common';
import type { Request, Response } from 'express';
import { Throttle } from '@nestjs/throttler';
import { CurrentUser } from '../common/decorators/current-user.decorator';
import { JwtAuthGuard } from '../common/guards/jwt-auth.guard';
import { User } from '../database/entities/user.entity';
import { buildUserResource } from '../users/user.resource';
import { AdminLoginDto } from './dto/admin-login.dto';
import { LoginDto } from './dto/login.dto';
import { PasswordResetDto } from './dto/password-reset.dto';
import { PasswordResetEmailDto } from './dto/password-reset-email.dto';
import { RegisterDto } from './dto/register.dto';
import { AuthService } from './auth.service';

@Controller()
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @Post('api/v1/user/register')
  async register(@Body() payload: RegisterDto) {
    return this.authService.register(payload);
  }

  @Post('api/v1/user/login')
  @HttpCode(200)
  async login(
    @Body() payload: LoginDto,
    @Req() request: Request,
    @Res({ passthrough: true }) response: Response,
  ) {
    const result = await this.authService.login(payload, request);
    response.setHeader('Authorization', `Bearer ${result.token}`);
    return result;
  }

  @Post('api/v1/user/logout')
  @UseGuards(JwtAuthGuard)
  @HttpCode(200)
  async logout(@CurrentUser() user: User) {
    return this.authService.logout(user);
  }

  @Get('api/v1/user/refresh')
  @UseGuards(JwtAuthGuard)
  async refresh(
    @CurrentUser() user: User,
    @Req() request: Request,
    @Res({ passthrough: true }) response: Response,
  ) {
    const result = await this.authService.refresh(user, request);
    response.setHeader('Authorization', `Bearer ${result.token}`);
    return result;
  }

  @Post('api/v1/user/password/email')
  @HttpCode(200)
  async sendResetEmail(@Body() payload: PasswordResetEmailDto) {
    return this.authService.sendPasswordResetEmail(payload);
  }

  @Post('api/v1/user/password/reset')
  @HttpCode(200)
  async resetPassword(@Body() payload: PasswordResetDto) {
    return this.authService.resetPassword(payload);
  }

  @Get('api/v1/user/email/verify/:id')
  @Throttle({ default: { limit: 6, ttl: 60_000 } })
  @Redirect(undefined, 302)
  async verifyEmail(
    @Param('id', ParseIntPipe) id: number,
    @Query('expires') expires: string,
    @Query('signature') signature: string,
  ) {
    const result = await this.authService.verifyEmail(id, expires, signature);
    return { url: result.redirectUrl };
  }

  @Get('api/v1/user/email/resend')
  @UseGuards(JwtAuthGuard)
  @Throttle({ default: { limit: 6, ttl: 60_000 } })
  async resendVerification(@CurrentUser() user: User) {
    return this.authService.resendVerificationEmail(user);
  }

  @Get('api/v1/user/email/success')
  async emailSuccess() {
    return { status: true, message: 'Email verification successful.' };
  }

  @Get('api/v1/user')
  @UseGuards(JwtAuthGuard)
  async getCurrentUser(@CurrentUser() user: User) {
    return buildUserResource(user, { includeProfile: true });
  }

  @Post('login')
  @HttpCode(200)
  async adminLogin(@Body() payload: AdminLoginDto, @Req() request: Request) {
    try {
      const admin = await this.authService.loginAdmin(payload.email, payload.password);
      request.session.adminId = admin.id;
      return { status: true, message: 'Admin login successful.' };
    } catch {
      throw new HttpException({ message: 'Admin login failed.' }, 501);
    }
  }

  @Get('logout')
  @Redirect('/admin', 302)
  async adminLogout(@Req() request: Request) {
    return new Promise((resolve, reject) => {
      request.session.destroy((error) => {
        if (error) {
          reject(error);
          return;
        }
        resolve({ url: '/admin' });
      });
    });
  }
}
