import { useEffect, useState } from 'react';
import './App.css';
import WordCloud from 'react-d3-cloud';

function App() {
  const [data, setData] = useState(null);

   const [words, setWords] = useState([
    { text: 'Hello', value: 10 },
    { text: 'World', value: 8 },
    { text: 'React', value: 15 },
    { text: 'Next.js', value: 12 },
    { text: 'Word Cloud', value: 30 },
    { text: 'D3.js', value: 18 }
   ]);
  
  useEffect(() => {
  const fetchData = async () => {
  const url = 'https://9641-223-190-85-154.ngrok.io/produce';
  const data = { message: 'bannana', type: 'fruit' };

  try {
    const response = await fetch(url, {
      method: 'POST',
      mode: 'cors',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    });

    const result = await response.json();
    console.log(result);
  } catch (error) {
    console.error(error);
  }
};
fetchData()
  }, []);

  const result = JSON.parse(data);
  // const words = Object.entries(result).map(([text, value]) => ({ text, value }));
 

  console.log("words===:>",words)

  return (
    <div className="text-[500px]">
      <h2>Word Cloud</h2>
      <WordCloud
        data={words}
        fontSizeMapper={word => Math.log2(word.value) * 5 + 10}
        rotate={word => word.value % 360}
      />
    </div>
  );
}

export default App;
