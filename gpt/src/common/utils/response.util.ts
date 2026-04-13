export function successResponse<T>(message: string, data: T, status = true) {
  return { status, message, data };
}

export function createdResponse<T>(status: string, message: string, data: T) {
  return { status, message, data };
}
