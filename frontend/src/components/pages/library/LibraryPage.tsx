import React, {useState, useEffect} from 'react';
import {message} from 'antd';
import BookList from './BookList';
import BookCreationForm from './BookCreationForm';
import {IBook} from '../../../types/bookTypes';
import api from '../../../api/api';
import axios from 'axios';
import {IAxiosError} from '../../../types/errorTypes';

const PAGE_SIZE = 10;

const LibraryPage: React.FC = () => {
  const [books, setBooks] = useState<IBook[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [hasMore, setHasMore] = useState<boolean>(true);
  const [page, setPage] = useState<number>(1);

  const fetchBooks = async (page: number) => {
    if (loading) return;
    setLoading(true);
    try {
      const response = await api.get(`/books`, {
        params: {page, page_size: PAGE_SIZE},
      });
      const newBooks: IBook[] = response.data.books;

      setBooks(prevBooks => {
        const existingBookIds = new Set(prevBooks.map(book => book.id));
        const filteredNewBooks = newBooks.filter(
          book => !existingBookIds.has(book.id),
        );
        return [...prevBooks, ...filteredNewBooks];
      });

      setHasMore(response.data.hasMore);
    } catch (error) {
      const axiosError = error as IAxiosError;
      if (axiosError.response) {
        message.error(
          'Failed to load books: ' +
            (axiosError.response.data.message || 'Network error'),
        );
      } else {
        message.error('Failed to load books: Unknown error');
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchBooks(page);
  }, [page]);

  const handleScroll = (event: React.UIEvent<HTMLElement>) => {
    const {scrollTop, clientHeight, scrollHeight} = event.currentTarget;
    if (scrollHeight - scrollTop <= clientHeight + 10 && hasMore && !loading) {
      setPage(prevPage => prevPage + 1);
    }
  };

  const handleCreateBook = async (title: string) => {
    try {
      const response = await api.post(`/books`, {title});
      const newBook: IBook = response.data;
      setBooks(prevBooks => [...prevBooks, newBook]);
      message.success('New book added!');
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        message.error(
          'Failed to add book: ' + error.response.data.message ||
            'Network error',
        );
      } else {
        message.error('Failed to add book: Unknown error');
      }
    }
  };

  const updateBookTitle = async (id: number, newTitle: string) => {
    try {
      await api.put(`/books/${id}`, {title: newTitle});
      setBooks(prevBooks => {
        return prevBooks.map(book =>
          book.id === id ? {...book, title: newTitle} : book,
        );
      });
      message.success('Book title updated successfully.');
    } catch (error) {
      const axiosError = error as IAxiosError;
      message.error(
        'Failed to update book title: ' +
          (axiosError.response?.data.message || 'Network error'),
      );
    }
  };

  const updateBookRating = async (id: number, rating: number) => {
    try {
      await api.put(`/books/${id}`, {rating});
      setBooks(prevBooks => {
        return prevBooks.map(book =>
          book.id === id ? {...book, rating} : book,
        );
      });
      message.success('Book rating updated successfully.');
    } catch (error) {
      const axiosError = error as IAxiosError;
      message.error(
        'Failed to update book rating: ' +
          (axiosError.response?.data.message || 'Network error'),
      );
    }
  };

  const handleDeleteBook = async (id: number) => {
    try {
      await api.delete(`/books/${id}`);
      setBooks(prevBooks => prevBooks.filter(book => book.id !== id));
      message.success('Book deleted!');
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        message.error(
          'Failed to delete book: ' + error.response.data.message ||
            'Network error',
        );
      } else {
        message.error('Failed to delete book: Unknown error');
      }
    }
  };

  return (
    <div
      style={{
        height: '100%',
        padding: '20px',
        fontSize: '18px',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'space-between',
      }}>
      <BookCreationForm onCreate={handleCreateBook} />
      <div
        style={{
          height: '500px',
          overflowY: 'auto',
          border: '1px solid #eaeaea',
          padding: '10px',
        }}
        onScroll={handleScroll}>
        <BookList
          books={books}
          loading={loading}
          onDelete={handleDeleteBook}
          onRatingChange={updateBookRating}
          onTitleChange={updateBookTitle}
        />
      </div>
    </div>
  );
};

export default LibraryPage;
