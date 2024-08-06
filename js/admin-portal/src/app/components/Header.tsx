import Link from 'next/link';

const Header = () => {
  return (
    <header className="bg-blue-500 p-4">
      <nav className="flex space-x-4">
        <Link href="/pages/page1" className="text-white">Page 1</Link>
        <Link href="/pages/page2" className="text-white">Page 2</Link>
        <Link href="/pages/page3" className="text-white">Page 3</Link>
        <Link href="/pages/index" className="text-white">Page 4</Link>
        <Link href="/" className="text-white">Home</Link>
      </nav>
    </header>
  );
};

export default Header;