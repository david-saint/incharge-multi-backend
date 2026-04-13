import {
  ArgumentsHost,
  Catch,
  ExceptionFilter,
  HttpException,
  UnauthorizedException,
  UnprocessableEntityException,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import type { Request, Response } from 'express';

@Catch()
export class HttpExceptionFilter implements ExceptionFilter {
  constructor(private readonly config: ConfigService) {}

  catch(exception: unknown, host: ArgumentsHost): void {
    const response = host.switchToHttp().getResponse<Response>();
    const request = host.switchToHttp().getRequest<Request>();

    if (exception instanceof Error) {
      const validationPayload = this.parseValidationPayload(exception);
      if (validationPayload) {
        response.status(422).json({ errors: validationPayload.errors });
        return;
      }
    }

    if (exception instanceof UnprocessableEntityException) {
      response.status(422).json(exception.getResponse());
      return;
    }

    if (exception instanceof UnauthorizedException) {
      response.status(401).json({ message: 'Permission Denied' });
      return;
    }

    if (exception instanceof HttpException) {
      response.status(exception.getStatus()).json(exception.getResponse());
      return;
    }

    const isProduction = this.config.get<string>('app.env') === 'production';
    const error = exception instanceof Error ? exception : new Error(String(exception));

    response.status(500).json({
      status: false,
      message: `Unhandled error while processing ${request.method} ${request.url}`,
      ...(isProduction
        ? {}
        : {
            error: error.message,
            trace: error.stack,
          }),
    });
  }

  private parseValidationPayload(exception: Error): {
    errors: Record<string, string[]>;
  } | null {
    try {
      const payload = JSON.parse(exception.message) as {
        __validation__?: boolean;
        errors?: Record<string, string[]>;
      };
      if (!payload.__validation__ || !payload.errors) {
        return null;
      }
      return { errors: payload.errors };
    } catch {
      return null;
    }
  }
}
