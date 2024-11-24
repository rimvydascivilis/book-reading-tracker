export interface IGoal {
  type: 'books' | 'pages';
  frequency: 'daily' | 'monthly';
  value: number;
}

export interface IGoalProgress {
  percentage: number;
  left: number;
}
