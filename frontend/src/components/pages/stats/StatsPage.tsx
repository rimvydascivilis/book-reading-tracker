import React, {useState, useEffect} from 'react';
import {DatePicker, Select, message, Row, Col, Typography, Card} from 'antd';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ReferenceLine,
} from 'recharts';
import api from '../../../api/api';
import {IAxiosError} from '../../../types/errorTypes';
import moment from 'moment';
import {IStatResponse} from '../../../types/statTypes';

const {Option} = Select;
const {MonthPicker} = DatePicker;
const {Title} = Typography;

const StatsPage: React.FC = () => {
  const [frequency, setFrequency] = useState<'monthly' | 'daily'>('monthly');
  const [year, setYear] = useState<number | null>(null);
  const [month, setMonth] = useState<number | null>(null);
  const [stats, setStats] = useState<IStatResponse | null>(null);

  useEffect(() => {
    if (!year) return;
    if (frequency === 'daily' && !month) return;

    const fetchStats = async () => {
      try {
        const endpoint =
          frequency === 'monthly'
            ? `/stats/monthly?year=${year}`
            : `/stats/daily?year=${year}&month=${month}`;
        const response = await api.get<IStatResponse>(endpoint);
        const fetchedStats = response.data;
        if (!fetchedStats.progress) fetchedStats.progress = [];

        if (frequency === 'monthly') {
          const months = Array.from({length: 12}, (_, i) => String(i + 1));
          fetchedStats.progress = months.map(month => {
            const existing = fetchedStats.progress.find(p => p.date === month);
            return existing || {date: month, pages: 0};
          });
        } else if (frequency === 'daily' && month) {
          const daysInMonth = new Date(year, month, 0).getDate();
          const days = Array.from({length: daysInMonth}, (_, i) =>
            String(i + 1),
          );
          fetchedStats.progress = days.map(day => {
            const existing = fetchedStats.progress.find(p => p.date === day);
            return existing || {date: day, pages: 0};
          });
        }

        setStats(fetchedStats);
      } catch (error) {
        const axiosError = error as IAxiosError;
        message.error(
          'Failed to fetch stats: ' + axiosError.response?.data.message,
        );
      }
    };

    fetchStats();
  }, [year, month, frequency]);

  const data =
    stats?.progress.map(item => ({
      date: item.date,
      pages: item.pages,
    })) || [];

  const maxPages = Math.max(...data.map(item => item.pages), 0);
  const yAxisMax =
    stats?.goal && stats.goal > maxPages ? stats.goal * 1.1 : undefined;

  return (
    <div style={{padding: 24}}>
      <Card>
        <Title level={2} style={{textAlign: 'center', marginBottom: 24}}>
          Reading Stats
        </Title>

        <Row gutter={16} justify="center" style={{marginBottom: 24}}>
          <Col>
            <Select
              value={frequency}
              onChange={setFrequency}
              style={{width: 120}}>
              <Option value="monthly">Monthly</Option>
              <Option value="daily">Daily</Option>
            </Select>
          </Col>

          <Col>
            <DatePicker
              picker="year"
              onChange={(date: moment.Moment | null) =>
                setYear(date ? date.year() : null)
              }
              placeholder="Select Year"
            />
          </Col>

          {frequency === 'daily' && (
            <Col>
              <MonthPicker
                onChange={(date: moment.Moment | null) =>
                  setMonth(date ? date.month() + 1 : null)
                }
                placeholder="Select Month"
                disabled={!year}
              />
            </Col>
          )}
        </Row>

        {stats ? (
          <BarChart
            style={{margin: 'auto'}}
            width={800}
            height={400}
            data={data}
            margin={{top: 20, right: 30, left: 20, bottom: 5}}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis domain={[0, yAxisMax || 'auto']} />
            <Tooltip />
            <Legend />
            <Bar dataKey="pages" fill="#8884d8" name="Pages Read" />
            {stats.goal && (
              <ReferenceLine
                y={stats.goal}
                label="Goal"
                stroke="red"
                strokeDasharray="3 3"
              />
            )}
          </BarChart>
        ) : (
          <Title level={4} style={{textAlign: 'center'}}>
            No Data Available
          </Title>
        )}
      </Card>
    </div>
  );
};

export default StatsPage;
