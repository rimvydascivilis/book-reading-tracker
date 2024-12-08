import React, {useState} from 'react';
import {Input, Modal} from 'antd';

const {TextArea} = Input;

interface NotesModalProps {
  visible: boolean;
  onClose: () => void;
  onCreate: (note: {page_number: number; content: string}) => void;
}

const NotesModal: React.FC<NotesModalProps> = ({
  visible,
  onClose,
  onCreate,
}) => {
  const [note, setNote] = useState<{page_number: number; content: string}>({
    page_number: 1,
    content: '',
  });

  const handleCreate = () => {
    onCreate(note);
    setNote({page_number: 1, content: ''});
  };

  return (
    <Modal
      title="Add Note"
      open={visible}
      onCancel={onClose}
      onOk={handleCreate}>
      <Input
        type="number"
        placeholder="Page Number"
        value={note.page_number}
        onChange={e => setNote({...note, page_number: Number(e.target.value)})}
      />
      <TextArea
        placeholder="Note Content"
        rows={4}
        value={note.content}
        onChange={e => setNote({...note, content: e.target.value})}
      />
    </Modal>
  );
};

export default NotesModal;
