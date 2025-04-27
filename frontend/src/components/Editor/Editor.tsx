import Header from '@editorjs/header'
import List from '@editorjs/list'
import Quote from '@editorjs/quote'
import Delimiter from '@editorjs/delimiter'
import InlineCode from '@editorjs/inline-code'
import Embed from "@editorjs/embed"
import EditorJS, { ToolConstructable } from '@editorjs/editorjs'
import Paragraph from "@editorjs/paragraph"
import UnderLine from "@editorjs/underline"
import { FC, useEffect, useRef } from 'react'
import { EditorFormProps } from './Editor.props'


const Editor: FC<EditorFormProps> = ({onChange}) => {
    const ref = useRef<EditorJS | null>(null);
    useEffect(() => {
        if (!ref.current) {
            const editor = new EditorJS({
                holder: "editorjs",
                onReady: () => {
                    console.log("Editor.js is ready");
                  },
                async onChange(api, event) {
                    const data = await api.saver.save()
                    onChange(data)
                },
                tools: {
                    paragraph: {
                        class: Paragraph as ToolConstructable,
                    },
                    underline: UnderLine,
                    header: {
                        class: Header as unknown as ToolConstructable,
                        config: {
                            placeholder: "Заголовок...",
                            levels: [1, 2, 3], // H1 и H2
                            defaultLevel: 1,
                        }
                    },
                    embed: {
                        class: Embed as ToolConstructable,
                        config: {
                          services: {
                            youtube: true,
                            vimeo: true,
                          },
                        },
                      },
                    list: {
                        class: List as unknown as ToolConstructable,
                        inlineToolbar: true,
                      },
                      quote: {
                        class: Quote,
                        inlineToolbar: true,
                      },
                    delimiter: Delimiter,
                    inlineCode: InlineCode
                  }
            })
        }
        return () => {
            if (ref.current && ref.current.destroy) {
                ref.current.destroy()
            }
        }
    }, [])
  return (
    <div id={"editorjs"} className="bg-base-lightBlue/30 rounded-lg max-w-3xl mx-auto my-8">
    </div>
  )
}
export default Editor;
// Загрузка изображений на сервер
// async function uploadToServer(file: File) {
//   const formData = new FormData()
//   formData.append('image', file)
  
//   const response = await fetch('/api/upload', {
//     method: 'POST',
//     body: formData
//   })
  
//   const data = await response.json()
//   return data.url
// }