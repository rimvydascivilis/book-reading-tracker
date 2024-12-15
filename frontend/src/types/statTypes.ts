export interface IProgressData {
  date: string;
  pages: number;
}

export interface IStatResponse {
  progress: IProgressData[];
  goal: number;
}
