export interface IReading {
  id: number;
  total_pages: number;
  link?: string;
}

export interface ICombinedReading {
  book_title: string;
  status: 'not started' | 'reading' | 'completed';
  progress: number;
  reading: IReading;
}

export interface IReadingFormValues {
  book_title: string;
  book_id: number;
  total_pages: number;
  link?: string;
}
