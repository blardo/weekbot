'use client'

import { useState, useEffect } from 'react';

const Sidebar = () => {
  const [view, setView] = useState('view1');
  const [collapsed, setCollapsed] = useState(false);
  const [pingData, setPingData] = useState<string | null>(null);

  useEffect(() => {
    const fetchPingData = async () => {
      try {
        console.log('Sending request to /api/ping');
        const response = await fetch('http://localhost:8080/ping');
        console.log('Received response:', response);
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const data = await response.text();
        setPingData(data);
      } catch (error) {
        console.error('Error fetching ping:', error);
      }
    };
  
    fetchPingData();
  }, [view]);

  return (
    <aside className={`bg-gray-200 p-4 ${collapsed ? 'w-16' : 'w-1/4'} transition-all duration-300`}>
      <button 
        onClick={() => setCollapsed(!collapsed)} 
        className="bg-blue-500 text-white p-2 mb-4"
      >
        {collapsed ? '>>' : '<<'}
      </button>
      <div className={`flex flex-col space-y-2 ${collapsed ? 'hidden' : 'block'}`}>
        <button onClick={() => setView('view1')} className="bg-blue-500 text-white p-2">View 1</button>
        <button onClick={() => setView('view2')} className="bg-blue-500 text-white p-2">View 2</button>
        <button onClick={() => setView('view3')} className="bg-blue-500 text-white p-2">View 3</button>
        <button onClick={() => setView('view4')} className="bg-blue-500 text-white p-2">View 4</button>
      </div>
      <div className="mt-4">
        {view === 'view1' && <div>{collapsed ? "" : "Content for View 1"}</div>}
        {view === 'view2' && <div>{collapsed ? "" : "Content for View 2"}</div>}
        {view === 'view3' && <div>{collapsed ? "" : "Content for View 3"}</div>}
        {view === 'view4' && <div>{collapsed ? "" : "Content for View 4"}</div>}
      </div>
      <div className="mt-4">
        {pingData && <div>Ping Data: {pingData}</div>}
      </div>
    </aside>
  );
};

export default Sidebar;