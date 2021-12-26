export interface IApiError {
  error?: { [key: string]: string }[] | string;
}

export default class ApiError {
  error?: { [key: string]: string }[] | string;

  constructor({ error }: IApiError) {
    this.error = error;
  }

  getFormattedStatus(): string | undefined {
    if (typeof this.error === "string") {
      return this.error;
    }
  }

  getError(): { [key: string]: string } {
    const errs = Object();

    if (typeof this.error === "object") {
      this.error.forEach((e) => {
        Object.entries(e).forEach(([key, value]) => (errs[key] = value));
      });
    }

    return errs;
  }
}
