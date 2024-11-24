import React, {useCallback} from 'react';
import {Table, Space, Progress, Button} from 'antd';
import {ICombinedReading} from '../../../types/readingTypes';
import {PlusOutlined} from '@ant-design/icons';

interface ReadingTableProps {
  readings: ICombinedReading[];
  onAddProgress: (id: number) => void;
  footer: React.ReactNode;
  loadMoreData: () => void;
  isLoading: boolean;
  hasMore: boolean;
}

const ReadingTable: React.FC<ReadingTableProps> = ({
  readings,
  onAddProgress,
  footer,
  loadMoreData,
  isLoading,
  hasMore,
}) => {
  const columns = [
    {
      title: 'Book Title',
      dataIndex: 'book_title',
      key: 'bookTitle',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
    },
    {
      title: 'Pages Read',
      key: 'pagesRead',
      render: (_: string, record: ICombinedReading) => (
        <Space>
          {record.progress} / {record.reading.total_pages}
        </Space>
      ),
    },
    {
      title: 'Progress',
      key: 'progress',
      render: (_: string, record: ICombinedReading) => (
        <Progress
          percent={Math.round(
            (record.progress / record.reading.total_pages) * 100,
          )}
          status="active"
          strokeColor="#1890ff"
        />
      ),
    },
    {
      title: 'Add Pages',
      key: 'addPages',

      render: (_: string, record: ICombinedReading) => (
        <Button onClick={() => onAddProgress(record.reading.id)}>
          <PlusOutlined />
        </Button>
      ),
    },
  ];

  const handleScroll = useCallback(
    (e: React.UIEvent<HTMLElement>) => {
      const {scrollHeight, scrollTop, clientHeight} = e.target as HTMLElement;
      if (
        scrollHeight - scrollTop <= clientHeight + 10 &&
        hasMore &&
        !isLoading
      ) {
        loadMoreData();
      }
    },
    [isLoading, loadMoreData],
  );

  return (
    <div style={{height: '100%'}} onScroll={handleScroll}>
      <Table
        columns={columns}
        dataSource={readings}
        rowKey="bookTitle"
        footer={() => footer}
        pagination={false}
        expandable={{
          expandedRowRender: (record: ICombinedReading) => (
            <a
              href={record.reading.link}
              target="_blank"
              rel="noopener noreferrer">
              Go to Book
            </a>
          ),
          rowExpandable: (record: ICombinedReading) =>
            record.reading.link ? true : false,
        }}
        scroll={{
          y: 250,
        }}
        loading={isLoading}
      />
    </div>
  );
};

export default ReadingTable;
