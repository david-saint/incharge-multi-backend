import { Controller, Delete, Get, Param, ParseIntPipe, Put, UseGuards } from '@nestjs/common';
import { AdminSessionGuard } from '../common/guards/admin-session.guard';
import { VerifiedAdminGuard } from '../common/guards/verified-admin.guard';
import { UsersService } from './users.service';

@Controller()
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  @Get('api/v1/user/users')
  @UseGuards(AdminSessionGuard)
  listUsersApi() {
    return this.usersService.listDeleted(false, '/api/v1/user/users');
  }

  @Get('api/v1/user/users/deletedUser')
  @UseGuards(AdminSessionGuard)
  listDeletedUsersApi() {
    return this.usersService.listDeleted(true, '/api/v1/user/users/deletedUser');
  }

  @Put('api/v1/user/users/update/:user_id')
  @UseGuards(AdminSessionGuard)
  restoreUserApi(@Param('user_id', ParseIntPipe) userId: number) {
    return this.usersService.restore(userId);
  }

  @Delete('api/v1/user/users/deleteUser/:user_id')
  @UseGuards(AdminSessionGuard)
  deleteUserApi(@Param('user_id', ParseIntPipe) userId: number) {
    return this.usersService.softDelete(userId);
  }

  @Get('getUsers')
  @UseGuards(VerifiedAdminGuard)
  listUsers() {
    return this.usersService.listDeleted(false);
  }

  @Get('getDeletedUsers')
  @UseGuards(VerifiedAdminGuard)
  listDeletedUsers() {
    return this.usersService.listDeleted(true);
  }

  @Put('updateUser/:user_id')
  @UseGuards(VerifiedAdminGuard)
  updateUser(@Param('user_id', ParseIntPipe) userId: number) {
    return this.usersService.restore(userId);
  }

  @Delete('deleteUser/:user_id')
  @UseGuards(VerifiedAdminGuard)
  deleteUser(@Param('user_id', ParseIntPipe) userId: number) {
    return this.usersService.softDelete(userId);
  }

  @Put('revertDeletedUser/:user_id')
  @UseGuards(VerifiedAdminGuard)
  revertDeletedUser(@Param('user_id', ParseIntPipe) userId: number) {
    return this.usersService.restore(userId);
  }
}
