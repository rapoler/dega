import React from 'react';
import EditorJS from '@editorjs/editorjs';
import Header from '@editorjs/header';
import List from '@editorjs/list';
import Paragraph from '@editorjs/paragraph';
import Quote from '@editorjs/quote';
import Table from '@editorjs/table';

function Editor({ value, onChange }) {
  React.useEffect(() => {
    new EditorJS({
      holder: 'editorjs',
      tools: {
        header: Header,
        list: List,
        paragraph: Paragraph,
        quote: Quote,
        table: Table,
      },
      onChange: (value) => value.saver.save().then((value) => onChange(value)),
      data: value,
    });
  }, []);

  return <div id="editorjs"></div>;
}

export default Editor;
