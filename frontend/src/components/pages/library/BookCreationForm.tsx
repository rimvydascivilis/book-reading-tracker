import React, {useState} from "react";
import {Input, Button, message} from "antd";
import { MAX_BOOK_TITLE_LENGTH } from "../../../constants";

interface BookCreationFormProps {
  onCreate: (title: string) => void;
}

const BookCreationForm: React.FC<BookCreationFormProps> = ({onCreate}) => {
  const [newBookTitle, setNewBookTitle] = useState<string>("");

  const handleCreateBook = () => {
    if (!newBookTitle.trim()) {
      message.error("Book title is required!");
      return;
    }
    onCreate(newBookTitle);
    setNewBookTitle("");
  };

  return (
    <div
      style={{
        marginBottom: "20px",
        display: "flex",
        justifyContent: "space-between",
      }}>
      <Input
        placeholder="Enter book title"
        value={newBookTitle}
        onChange={e => setNewBookTitle(e.target.value)}
        maxLength={MAX_BOOK_TITLE_LENGTH}
        style={{width: "70%", marginRight: "10px"}}
      />
      <Button type="primary" onClick={handleCreateBook}>
        Create New Book
      </Button>
    </div>
  );
};

export default BookCreationForm;
