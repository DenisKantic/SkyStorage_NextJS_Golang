"use client"
import React, {useState, useEffect} from 'react'
import FormUpload from './FormUpload'
import { getFiles, File } from '../../../actions/getFiles';
import { BsThreeDotsVertical } from "react-icons/bs";
import { FaRegTrashAlt } from "react-icons/fa";
import { MdOutlineFileDownload } from "react-icons/md";
import { BiRename } from "react-icons/bi";
import { deleteFile } from '../../../actions/deleteFile';

 
const Dashboard = () => {

  const [files, setFiles] = useState<File[]>([]);
  const fetchFiles = async () => {
    const data = await getFiles();
    setFiles(data);
  };
  
  // Fetch files when component mounts
  useEffect(() => {
    fetchFiles();
  }, []);
  


function bytesToGB(bytes:number): string{
  if (bytes === 0) return '0 B';

  const MB = bytes / (1024 * 1024);
  
  if (MB < 1000) {
    return `${MB.toFixed(2)} MB`;
  } else {
    const GB = MB / 1024;
    return `${GB.toFixed(2)} GB`;
  }
}

const handleDelete = async (id:number)=>{
  const response = await deleteFile(id);

  alert("Delete okay")
  fetchFiles();

}

  
  return (
    <div className='h-[100svh] w-full bg-white px-10'>
      <div className='grid grid-cols-2 gap-10 pt-20'>
        <FormUpload />
      </div>

      <div className='pt-10'>
        <p className='text-2xl text-black'>Files</p>
        <div className="overflow-x-auto text-black  w-full min-h-[50vh]">
         <table className="table">
           {/* head */}
           <thead className='font-bold text-black'>
             <tr>
               <th>ID</th>
               <th>Name</th>
               <th>File type</th>
               <th>File size</th>
               <th>Uploaded_at</th>
               <th>Options</th>
             </tr>
           </thead>
           <tbody>
             {/* row 1 */}
        {files.map((file) => (
       
             <tr key={file.id}>
               <th>{file.id}</th>
               <td>{file.filename}</td>
               <td>{file.filetype}</td>
               <td>{bytesToGB(file.filesize)}</td>
               <td>{file.uploaded_at.toDateString()}</td>
               <td>
               <div className="dropdown dropdown-left w-full flex items-center justify-center">
                  <div tabIndex={0} role="button" className="text-xl"><BsThreeDotsVertical /></div>
                  <ul tabIndex={0} className="dropdown-content menu shadow-2xl bg-white rounded-box z-[1] w-52 p-2 ">
                    <li className='pb-2'><a><MdOutlineFileDownload /> Download</a></li>
                    <li className='pb-2'><a><BiRename /> Rename</a></li>
                    <li onClick={()=>handleDelete(file.id)} className='bg-red-600 rounded-xl text-white font-bold'><a><FaRegTrashAlt /> Delete</a></li>
                  </ul>
                </div>
               </td>
             </tr>
        ))}
        </tbody>
      </table>
    </div>
      </div>
    </div>
  )
}

export default Dashboard
