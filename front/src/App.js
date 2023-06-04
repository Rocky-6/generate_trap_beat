import React, { useState } from 'react';
import axios from 'axios';

const App = () => {
  const [selectedKey, setSelectedKey] = useState('');

  const handleKeyChange = (event) => {
    setSelectedKey(event.target.value);
  };

  const handleDownload = async () => {
    try {
      const response = await axios.get("http://localhost:8080/download?key=" + selectedKey, {
        responseType: 'blob',
      });

      const url = URL.createObjectURL(response.data);

      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'midi-files.zip');
      link.click();
    } catch (error) {
      console.error('Error:', error);
    }
  };

  return (
    <div>
      <h1>自動 Trap Beat 生成</h1>
      <select value={selectedKey} onChange={handleKeyChange}>
        <option value="">キーを選択してください</option>
        <option value="C">C</option>
        <option value="C#">C#</option>
        <option value="D">D</option>
        <option value="D#">D#</option>
        <option value="E">E</option>
        <option value="F">F</option>
        <option value="F#">F#</option>
        <option value="G">G</option>
        <option value="G#">G#</option>
        <option value="A">A</option>
        <option value="A#">A#</option>
        <option value="B">B</option>
        {/* 他のキーのオプションを追加 */}
      </select>
      <button onClick={handleDownload}>ダウンロード</button>
    </div>
  );
};

export default App;
