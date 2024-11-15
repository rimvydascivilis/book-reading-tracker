export interface IGoal {
  type: 'books' | 'pages';
  frequency: 'daily' | 'monthly';
  value: number;
}
