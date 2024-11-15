export interface IAxiosError {
  response?: {
    data: {
      message: string;
    };
    status: number;
  };
}
