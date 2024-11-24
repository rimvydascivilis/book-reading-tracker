import React from 'react';
import {Form, InputNumber, Button, DatePicker, Modal} from 'antd';

interface AddProgressModalProps {
  open: boolean;
  readingId: number | null;
  bookTitle: string;
  onAddProgress: (
    readingId: number,
    pagesRead: number,
    readingDate: Date,
  ) => void;
  onCancel: () => void;
}

const AddProgressModal: React.FC<AddProgressModalProps> = ({
  open,
  readingId,
  bookTitle,
  onAddProgress,
  onCancel,
}) => {
  const [form] = Form.useForm();

  const handleAddProgress = (values: {
    pagesRead: number;
    readingDate: Date;
  }) => {
    if (!readingId) return;
    onAddProgress(readingId, values.pagesRead, values.readingDate);
    form.resetFields();
  };

  return (
    <Modal
      title={`Add Progress for "${bookTitle}"`}
      open={open}
      onCancel={() => {
        form.resetFields();
        onCancel();
      }}
      footer={null}>
      <Form form={form} onFinish={handleAddProgress} layout="vertical">
        <Form.Item
          label="Pages Read"
          name="pagesRead"
          rules={[{required: true, message: 'Please enter the pages read!'}]}>
          <InputNumber
            min={1}
            placeholder="Enter pages read"
            style={{width: '100%'}}
          />
        </Form.Item>
        <Form.Item
          label="Reading Date"
          name="readingDate"
          rules={[{required: true, message: 'Please select a date!'}]}>
          <DatePicker style={{width: '100%'}} />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" style={{width: '100%'}}>
            Add Progress
          </Button>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default AddProgressModal;
