import React from 'react';
import {Modal, Button, Typography, Input} from 'antd';
import {ListItem} from '../../../types/listTypes';

const {TextArea} = Input;

interface ExportListModalProps {
  visible: boolean;
  onClose: () => void;
  listItems: ListItem[];
  listName: string;
}

const ExportListModal: React.FC<ExportListModalProps> = ({
  visible,
  onClose,
  listItems,
  listName,
}) => {
  const generateReadableList = () => {
    return listItems
      .map((item, index) => `${index + 1}. ${item.book_name}`)
      .join('\n');
  };

  const handleCopyToClipboard = () => {
    navigator.clipboard.writeText(generateReadableList());
    Modal.success({
      content: 'List copied to clipboard!',
    });
  };

  const handleDownloadFile = () => {
    const blob = new Blob([generateReadableList()], {
      type: 'text/plain;charset=utf-8',
    });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `${listName}.txt`;
    link.click();
    URL.revokeObjectURL(link.href);
  };

  return (
    <Modal
      visible={visible}
      title={`Export List: ${listName}`}
      onCancel={onClose}
      footer={[
        <Button key="copy" onClick={handleCopyToClipboard}>
          Copy to Clipboard
        </Button>,
        <Button key="download" type="primary" onClick={handleDownloadFile}>
          Download as File
        </Button>,
        <Button key="cancel" onClick={onClose}>
          Close
        </Button>,
      ]}>
      <Typography.Text>Here is your list in a readable format:</Typography.Text>
      <TextArea
        value={generateReadableList()}
        readOnly
        rows={8}
        style={{marginTop: '12px'}}
      />
    </Modal>
  );
};

export default ExportListModal;
