import React, { useState, useEffect } from 'react';
import { Layout, Form, InputNumber, Button, Typography, Select, message } from 'antd';
import { IGoal } from '../../../types/goalTypes';
import { IAxiosError } from '../../../types/errorTypes';
import api from '../../../api/api';
import Loading from '../../common/Loading';

const { Title, Text } = Typography;
const { Option } = Select;

const GoalPage: React.FC = () => {
  const [goal, setGoal] = useState<IGoal | null>(null);

  useEffect(() => {
    fetchGoal();
  }, []);

  const fetchGoal = async () => {
    try {
      const response = await api.get('/goal');
      const data = response.data as IGoal;
      setGoal(data);
    } catch (error) {
      const axiosError = error as IAxiosError
      if (axiosError.response?.status === 404) {
        // Goal not found, set default goal
        setGoal({ type: 'books', frequency: 'daily', value: 0 });
        return;
      }

      message.error(
        'Failed to fetch current goal: ' +
          (axiosError.response?.data.message || 'Network error'),
      );
    }
  };

  const updateGoal = async (values: IGoal) => {
    try {
      const response = await api.put('/goal', values);

      if (!response.data) {
        throw new Error('Invalid response');
      }

      setGoal(values);
      message.success('Goal updated successfully');
    } catch (error) {
      const axiosError = error as IAxiosError;
      message.error('Failed to update goal: ' + (axiosError.response?.data.message || 'Network error'));
    }
  };

  const onFinish = (values: { type: IGoal['type']; frequency: IGoal['frequency']; value: number }) => {
    updateGoal({ type: values.type, frequency: values.frequency, value: values.value });
  };

  if (goal === null) {
    return <Loading />
  }

  return (
    <Layout
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100%',
        background: '#fff',
        padding: '24px',
      }}
    >
      <div style={{ maxWidth: '400px', width: '100%' }}>
        <Title level={3} style={{ textAlign: 'center' }}>Change Your Reading Goal</Title>
        <Form
          layout="vertical"
          onFinish={onFinish}
          initialValues={goal}
          style={{ marginTop: '24px' }}
        >
          <Form.Item
            label="Goal Type"
            name="type"
            rules={[{ required: true, message: 'Please select your goal type!' }]}
          >
            <Select>
              <Option value="books">Books</Option>
              <Option value="pages">Pages</Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="Goal Frequency"
            name="frequency"
            rules={[{ required: true, message: 'Please select your goal frequency!' }]}
          >
            <Select>
              <Option value="daily">Daily</Option>
              <Option value="monthly">Monthly</Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="Goal Value"
            name="value"
            rules={[{ required: true, message: 'Please enter your goal value!' }]}
          >
            <InputNumber min={1} placeholder="Enter goal value" style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              Set Goal
            </Button>
          </Form.Item>
        </Form>

        {goal.value > 0 && (
          <Text style={{ display: 'block', marginTop: '24px', textAlign: 'center' }}>
            Current Goal: <strong>{goal.value} {goal.type} ({goal.frequency})</strong>
          </Text>
        )}
      </div>
    </Layout>
  );
};

export default GoalPage;
