import React from 'react'
import Image from 'next/image'
import userImage from "../../../public/images/user.png"

const Navbar = () => {
  return (
    <div className="navbar bg-gray-200 fixed text-black h-[6svh]">
    <div className="flex-1">
      <a className="btn btn-ghost text-xl">SkyStorage</a>
    </div>
    <div className="flex-none gap-2">
      <div className="dropdown dropdown-end -">
        <div tabIndex={0} role="button" className="btn btn-ghost btn-circle avatar">
          <div className="w-10 rounded-full">
            <Image
              alt="Tailwind CSS Navbar component"
              src={userImage}
              height={50}
              width={50} />
          </div>
        </div>
        <ul
          tabIndex={0}
          className="menu menu-sm dropdown-content bg-white rounded-box z-[1] mt-3 w-52 p-2 shadow">
          <li>
            <a className="justify-between">
              Profile
              <span className="badge">New</span>
            </a>
          </li>
          <li><a>Settings</a></li>
          <li><a>Logout</a></li>
        </ul>
      </div>
    </div>
  </div>
  )
}

export default Navbar