//layout.js is for the components that are visible on every page - header? 
//information to see in the browsers tab and google search etc
export const metadata = {
  title: "Social Network",
  description: "Connect with friends",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body className="bg-white text-black">
        <nav className="flex flex-row space-x-4 p-4 bg-gray-100">
          <a href="/" className="text-blue-600 hover:underline">Main page</a>
          <a href="/login" className="text-blue-600 hover:underline">Login</a>
          <a href="/register" className="text-blue-600 hover:underline">Register</a>
        </nav>
        <main className="p-5">
          {children}
        </main>
      </body>
    </html>
  );
}

