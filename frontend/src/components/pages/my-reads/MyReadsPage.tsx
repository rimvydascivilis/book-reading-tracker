import React, {useState, useEffect} from 'react';
import {Button, message} from 'antd';
import {
  ICombinedReading,
  IReading,
  IReadingFormValues,
} from '../../../types/readingTypes';
import AddReadingModal from './AddReadingModal';
import AddProgressModal from './AddProgressModal';
import ReadingTable from './ReadingTable';
import {PlusOutlined} from '@ant-design/icons';
import api from '../../../api/api';
import GoalCard from './GoalCard';
import {IAxiosError} from '../../../types/errorTypes';
import {IGoalProgress} from '../../../types/goalTypes';

const MyReadsPage: React.FC = () => {
  const [readings, setReadings] = useState<ICombinedReading[]>([]);
  const [goalProgress, setGoalProgress] = useState<IGoalProgress>({
    percentage: 0,
    left: 0,
  });
  const [isReadingFormVisible, setIsReadingFormVisible] =
    useState<boolean>(false);
  const [selectedReadingId, setSelectedReadingId] = useState<number | null>(
    null,
  );
  const [isProgressModalVisible, setIsProgressModalVisible] =
    useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [hasMore, setHasMore] = useState<boolean>(true);
  const [page, setPage] = useState<number>(1);

  const fetchReadings = async (page: number) => {
    try {
      setIsLoading(true);
      const {data} = await api.get(`/readings?page=${page}`);
      const newReadings = data.readings as ICombinedReading[];
      const hasMore = data.hasMore as boolean;
      setReadings(prevReadings => [
        ...prevReadings,
        ...newReadings.filter(
          reading =>
            !prevReadings.some(
              prevReading => prevReading.reading.id === reading.reading.id,
            ),
        ),
      ]);
      setIsLoading(false);
      setHasMore(hasMore);
    } catch (error) {
      const axiosError = error as IAxiosError;
      message.error(
        'Failed to fetch readings: ' +
          (axiosError.response?.data.message || 'Network error'),
      );
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchReadings(page);
  }, [page]);

  const fetchGoalProgress = async () => {
    try {
      const {data} = await api.get('/goal/progress');
      setGoalProgress(data);
    } catch (error) {
      const axiosError = error as IAxiosError;
      if (axiosError.response?.status === 404) {
        message.info('Set a goal to track your reading progress');
      } else {
        message.error(
          'Failed to fetch goal progress: ' +
            (axiosError.response?.data.message || 'Network error'),
        );
      }
    }
  };

  useEffect(() => {
    fetchGoalProgress();
  }, []);

  const handleAddReading = async (newReading: IReadingFormValues) => {
    try {
      const resp = await api.post('/readings', {
        book_id: newReading.book_id,
        total_pages: newReading.total_pages,
        link: newReading.link,
      });
      const reading = resp.data as IReading;
      setReadings(prevReadings => [
        ...prevReadings,
        {
          book_title: newReading.book_title,
          status: 'not started',
          progress: 0,
          reading,
        },
      ]);
      setIsReadingFormVisible(false);
    } catch (error) {
      const axiosError = error as IAxiosError;
      message.error(
        'Failed to add reading: ' + axiosError.response?.data.message ||
          'unknown error',
      );
    }
  };

  const handleAddProgress = async (
    readingId: number,
    pagesRead: number,
    readingDate: Date,
  ) => {
    try {
      await api.post(`/progress/${readingId}`, {
        pages: pagesRead,
        date: readingDate,
      });
      const newReadings = readings.map(reading => {
        if (reading.reading.id === readingId) {
          return {
            ...reading,
            progress: reading.progress + pagesRead,
          };
        }
        return reading;
      });
      fetchGoalProgress(); // Update goal progress
      setReadings(newReadings);
      setIsProgressModalVisible(false);
    } catch (error) {
      const axiosError = error as IAxiosError;
      message.error(
        'Failed to add progress: ' + axiosError.response?.data.message ||
          'unknown error',
      );
    }
  };

  const loadMoreData = () => {
    if (isLoading) return;
    setPage(prevPage => prevPage + 1);
  };

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        height: '100%',
        padding: '2vh',
      }}>
      <GoalCard progress={goalProgress} />
      <ReadingTable
        readings={readings}
        onAddProgress={id => {
          setSelectedReadingId(id);
          setIsProgressModalVisible(true);
        }}
        footer={
          <Button
            type="primary"
            icon={<PlusOutlined />}
            size="large"
            onClick={() => setIsReadingFormVisible(true)}>
            Add New Reading
          </Button>
        }
        loadMoreData={loadMoreData}
        isLoading={isLoading}
        hasMore={hasMore}
      />

      <AddReadingModal
        open={isReadingFormVisible}
        onAddReading={handleAddReading}
        onCancel={() => setIsReadingFormVisible(false)}
      />

      <AddProgressModal
        open={isProgressModalVisible}
        readingId={selectedReadingId}
        bookTitle={
          readings.find(reading => reading.reading.id === selectedReadingId)
            ?.book_title || ''
        }
        onAddProgress={handleAddProgress}
        onCancel={() => setIsProgressModalVisible(false)}
      />
    </div>
  );
};

export default MyReadsPage;
