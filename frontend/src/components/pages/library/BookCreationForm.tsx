import React from 'react';
import {Form, Input, Button, message} from 'antd';
import {MAX_BOOK_TITLE_LENGTH} from '../../../constants';

interface BookCreationFormProps {
  onCreate: (title: string) => void;
}

const BookCreationForm: React.FC<BookCreationFormProps> = ({onCreate}) => {
  const [form] = Form.useForm();

  const handleCreateBook = (values: {title: string}) => {
    const {title} = values;
    if (!title.trim()) {
      message.error('Book title is required!');
      return;
    }
    onCreate(title);
    form.resetFields();
  };

  return (
    <Form
      form={form}
      layout="inline"
      onFinish={handleCreateBook}
      style={{
        display: 'flex',
        justifyContent: 'space-between',
        width: '100%',
      }}>
      <Form.Item
        name="title"
        rules={[
          {required: true, message: 'Book title is required!'},
          {
            max: MAX_BOOK_TITLE_LENGTH,
            message: `Title cannot exceed ${MAX_BOOK_TITLE_LENGTH} characters!`,
          },
        ]}
        style={{
          flex: 70,
          marginRight: '10px',
        }}>
        <Input
          placeholder="Enter book title"
          maxLength={MAX_BOOK_TITLE_LENGTH}
        />
      </Form.Item>
      <div style={{flex: 5}} />
      <Form.Item
        style={{
          flex: 25,
          maxWidth: '30%',
        }}>
        <Button type="primary" htmlType="submit" block>
          Create New Book
        </Button>
      </Form.Item>
    </Form>
  );
};

export default BookCreationForm;
