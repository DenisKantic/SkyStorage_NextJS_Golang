"use client"
import React from 'react'
import { uploadFile } from '../../../actions/uploadFile';
import { useRouter } from 'next/navigation';

export default function FormUpload() {

    const router = useRouter();

    const handleSubmit = async (e:React.FormEvent<HTMLFormElement>) =>{
        e.preventDefault();

        const form = e.currentTarget;
        const formData = new FormData(form);

        const data =await uploadFile(formData)
        const {path} = data;
        console.log("RESPONSE DATA", data, "PATH", path)
        alert(`Uploaded is complete.`)
        router.refresh();
    }

  return (
    <>
     <form onSubmit={handleSubmit} className='flex gap-10 '>
            <input type="file" name="file" className='w-full btn flex items-center justify-center btn-primary' placeholder='insert file' /> <br />
             <button type='submit' className='btn btn-info w-[50%]'>Upload</button>
        </form>
        </>
  )
}
