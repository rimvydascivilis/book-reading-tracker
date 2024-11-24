import React, {useEffect, useState} from 'react';
import {Card, Progress, message} from 'antd';
import api from '../../../api/api';
import {IAxiosError} from '../../../types/errorTypes';
import {IGoal, IGoalProgress} from '../../../types/goalTypes';

interface GoalCardProps {
  progress: IGoalProgress;
}

const GoalCard: React.FC<GoalCardProps> = ({progress}) => {
  const [goal, setGoal] = useState<IGoal | null>(null);

  const fetchGoal = async () => {
    try {
      const {data} = await api.get('/goal');
      setGoal(data);
    } catch (error) {
      const axiosError = error as IAxiosError;
      if (axiosError.response?.status === 404) {
        message.info('Set a goal to track your reading progress');
      } else {
        message.error(
          'Failed to fetch goal: ' +
            (axiosError.response?.data.message || 'Network error'),
        );
      }
    }
  };

  useEffect(() => {
    fetchGoal();
  }, []);

  if (!goal) {
    return null;
  }

  return (
    <Card style={{marginBottom: '16px', textAlign: 'center'}}>
      <h2>
        {goal.value} {goal.type} per{' '}
        {goal.frequency === 'daily' ? 'day' : 'month'}
      </h2>
      <Progress
        percent={progress.percentage}
        status={progress.left === 0 ? 'success' : 'active'}
      />

      {progress.left === 0 ? (
        <p>Congratulations! You have reached your goal</p>
      ) : (
        <p>
          You have {progress.left} {goal.type} left to read to reach your goal{' '}
          {goal.frequency === 'daily' ? 'today' : 'this month'}
        </p>
      )}
    </Card>
  );
};

export default GoalCard;
