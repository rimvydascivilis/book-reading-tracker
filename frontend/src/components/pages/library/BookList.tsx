import React from 'react';
import { List, Spin } from 'antd';
import BookItem from './BookItem';
import { IBook } from '../../../types/bookTypes';

interface BookListProps {
  books: IBook[];
  loading: boolean;
  onDelete: (id: number) => void;
  onRatingChange: (id: number, rating: number) => void;
  onTitleChange: (id: number, title: string) => void;
}

const BookList: React.FC<BookListProps> = ({ books, loading, onDelete, onRatingChange, onTitleChange }) => {
  return (
    <>
      <List
        dataSource={books}
        renderItem={
          (book) => <BookItem key={book.id} book={book}
        onDelete={onDelete} onRatingChange={onRatingChange} 
        onTitleChange={onTitleChange}
        />}
      />
      {loading && <Spin style={{ display: 'block', margin: '20px auto' }} />}
    </>
  );
};

export default BookList;
