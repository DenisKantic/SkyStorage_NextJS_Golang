"use client"
import React, {useState, useEffect} from 'react'
import FormUpload from './FormUpload'
import { getFiles, File } from '../../../actions/getFiles';
import { BsThreeDotsVertical } from "react-icons/bs";
import { FaRegTrashAlt } from "react-icons/fa";
import { MdOutlineFileDownload } from "react-icons/md";
import { BiRename } from "react-icons/bi";
import { deleteFile } from '../../../actions/deleteFile';
import { GrDocumentPdf } from "react-icons/gr";
import { TbFileTypeDocx } from "react-icons/tb";
import { BsFiletypeXls } from "react-icons/bs";
import { GrDocumentSound } from "react-icons/gr";




 
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

const getFileExtension = (filename:string) =>{
  
  let parts = filename.split('.');

  // Handle cases where the extension might be like 'application/pdf'
  if (parts.length > 1) {
    let lastPart = parts.pop()?.toLowerCase(); // Get the last part
    if (lastPart?.includes('/')) {
      // If the last part contains '/', split by '/' and take the first part
      return lastPart.split('/')[0];
    } else {
      return lastPart;
    }
  } else {
    return ''; // No extension found
  }
}
  
const getFileIcon = (filename:string) =>{
  
  const extension = getFileExtension(filename)

  switch(extension){
    case 'pdf':
      return <div className='text-xl text-red-800'><GrDocumentPdf /></div>
    case 'docx' :
    case 'doc' :
    case 'odt' :
      return <div className='text-xl text-blue-600'><TbFileTypeDocx /></div>
    case 'xls':
      case 'xlsx':
        case 'ods': // openDocument SpreadSheet
        case 'gsheet' : // google sheets
        case 'et' : //wps office spreadsheets
        return <div className="text-xl text-green-600"><BsFiletypeXls /></div>
    case 'mp3':
      return <div className="text-xl text-black"><GrDocumentSound /></div>

    default:
      return <div className='text-xl text-red-600'><BsThreeDotsVertical /></div>
  }
}

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

const handleDownloadFile = async (id: number, filename:string) => {
  try {
    const response = await fetch(`http://localhost:8080/download?id=${id}`, {
      method: 'GET',
    });

    if (!response.ok) {
      throw new Error('Failed to download file');
    }

    const fileExt = response.headers.get('X-File-Extension')
    console.log("TEST", fileExt)

    const blob = await response.blob();
    // Create download link
    const downloadUrl = window.URL.createObjectURL(new Blob([blob]));
    const a = document.createElement('a');
    a.href = downloadUrl;
    a.download = filename;
    a.style.display = 'none';
  

    // Append the anchor element to the body, click it, and remove it
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(downloadUrl);
    document.body.removeChild(a);
  } catch (error) {
    console.error('Error downloading file:', error);
  }
};



const downloadFiles = async (id:number, filetype:string)=>{
   
  try{
    await handleDownloadFile(id, filetype);
    fetchFiles();
  } catch(error){
    console.log("error downloading file", error)
  }
}

  
  return (
    <div className='h-[100svh] w-full bg-white px-10'>
        <p className='text-3xl text-black text-center pt-20'>Welcome to SkyStorage</p>
      <div className='grid grid-cols-2 gap-10 pt-10'>
        <FormUpload />
      </div>

      <div className='pt-10'>
        <h1 className='text-black text-2xl pb-5'>Storage {"(74% full)"}:</h1>
        <p className='text-black'>3.5GB of 15GB used</p>
        <progress className="progress progress-primary w-56 text-black" value={74} max="100"></progress>

        <p className='text-2xl text-black'>Files</p>
        <div className="overflow-x-auto text-black  w-full min-h-[50vh]">
         <table className="table">
           {/* head */}
           <thead className='font-bold text-black'>
             <tr>
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
               <td className='flex items-center'><span className='mr-5'>{getFileIcon(file.filename)}</span>{file.filename.substring(0,30)}</td>
               <td>{getFileExtension(file.filename)}</td>
               <td>{bytesToGB(file.filesize)}</td>
               <td>{file.uploaded_at.toDateString()}</td>
               <td>
               <div className="dropdown dropdown-left w-full flex items-center justify-center">
                  <div tabIndex={0} role="button" className="text-xl"><BsThreeDotsVertical /></div>
                  <ul tabIndex={0} className="dropdown-content menu shadow-2xl bg-white rounded-box z-[1] w-52 p-2 ">
                    <li onClick={()=>downloadFiles(file.id, file.filename)} className='pb-2'><a><MdOutlineFileDownload /> Download</a></li>
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
