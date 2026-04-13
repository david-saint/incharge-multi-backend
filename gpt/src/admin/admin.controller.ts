import { join } from 'node:path';
import {
  Body,
  Controller,
  Get,
  Param,
  ParseIntPipe,
  Post,
  Put,
  Req,
  Res,
  UnauthorizedException,
  UnprocessableEntityException,
  UseGuards,
} from '@nestjs/common';
import type { Request, Response } from 'express';
import { CreateAdminDto } from '../auth/dto/create-admin.dto';
import { VerifiedAdminGuard } from '../common/guards/verified-admin.guard';
import { AdminService } from './admin.service';
import { UpdateAdminDto } from './dto/update-admin.dto';

@Controller()
export class AdminController {
  constructor(private readonly adminService: AdminService) {}

  @Get()
  async root(@Res() response: Response) {
    const hasSuperAdmin = await this.adminService.hasSuperAdmin();
    if (hasSuperAdmin) {
      return response.redirect('/loginView');
    }
    return response.sendFile(join(process.cwd(), 'public', 'register-super-admin.html'));
  }

  @Get('admin')
  async adminRoot(@Res() response: Response) {
    const hasSuperAdmin = await this.adminService.hasSuperAdmin();
    if (hasSuperAdmin) {
      return response.redirect('/loginView');
    }
    return response.sendFile(join(process.cwd(), 'public', 'register-super-admin.html'));
  }

  @Get('loginView')
  loginView(@Res() response: Response) {
    return response.sendFile(join(process.cwd(), 'public', 'admin-login.html'));
  }

  @Get('privacy')
  privacyView(@Res() response: Response) {
    return response.sendFile(join(process.cwd(), 'public', 'privacy.html'));
  }

  @Get('panel')
  @UseGuards(VerifiedAdminGuard)
  panelView(@Res() response: Response) {
    return response.sendFile(join(process.cwd(), 'public', 'admin-panel.html'));
  }

  @Post('admin')
  async createAdmin(@Body() payload: CreateAdminDto, @Req() request: Request) {
    const hadSuperAdmin = await this.adminService.hasSuperAdmin();
    if (hadSuperAdmin && !request.session.adminId) {
      throw new UnauthorizedException();
    }
    if (hadSuperAdmin && request.session.adminId) {
      await this.adminService.getVerifiedAdmin(request.session.adminId);
    }
    if (!hadSuperAdmin && payload.userType !== 'Super') {
      throw new UnprocessableEntityException({
        errors: { userType: ['The first admin must be a Super admin.'] },
      });
    }

    const normalizedPayload =
      !hadSuperAdmin && payload.userType === 'Super'
        ? { ...payload, verified: 'Y' as const }
        : payload;
    const admin = await this.adminService.createAdmin(normalizedPayload);
    if (!hadSuperAdmin && normalizedPayload.userType === 'Super') {
      request.session.adminId = admin.id;
    }
    return admin;
  }

  @Get('allAdmins')
  @UseGuards(VerifiedAdminGuard)
  listAdmins() {
    return this.adminService.listAdmins();
  }

  @Get('getAdminDet')
  @UseGuards(VerifiedAdminGuard)
  getAdminDet(@Req() request: Request) {
    return this.adminService.getAdminDetails(request.session.adminId!);
  }

  @Put('updateAdmin/:admin_id')
  @UseGuards(VerifiedAdminGuard)
  updateAdmin(
    @Param('admin_id', ParseIntPipe) adminId: number,
    @Body() payload: UpdateAdminDto,
  ) {
    return this.adminService.updateAdmin(adminId, payload);
  }

  @Get('getContraceptionReason')
  @UseGuards(VerifiedAdminGuard)
  async getContraceptionReason() {
    const data = await this.adminService.referenceData();
    return data.reasons;
  }

  @Get('getEducationalLevels')
  @UseGuards(VerifiedAdminGuard)
  async getEducationalLevels() {
    const data = await this.adminService.referenceData();
    return data.educationLevels;
  }
}
