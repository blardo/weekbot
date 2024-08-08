import Link from 'next/link';

const Footer = () => {
  return (
    <footer className="bg-blue-500 p-4 mt-auto">
      <nav className="flex space-x-4">
      <Link href="/page1" className="text-white">Page 1</Link>
      <Link href="/page2" className="text-white">Page 2</Link>
      <Link href="/page3" className="text-white">Page 3</Link>
      <Link href="/page4" className="text-white">Page 4</Link>
      </nav>
    </footer>
  );
};

export default Footer;