import './globals.css';
import Header from './components/Header';
import Sidebar from './components/Sidebar';
import Footer from './components/Footer';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="flex flex-col min-h-screen">
        <Header />
        <div className="flex flex-1">
          <Sidebar />
          <main className="flex-1 p-4">
            {children}
          </main>
        </div>
        <Footer />
      </body>
    </html>
  );
}