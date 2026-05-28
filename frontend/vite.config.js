import { defineConfig } from 'vite'           //从vite库导入defineConfig函数
import vue from '@vitejs/plugin-vue'         //装vue翻译器
import path from 'path'                      //从Node.js里引用path函数的处理

export default defineConfig({                //外部可引用defineConfig这个函数 下面是函数内容
  plugins: [vue()],                         //告诉这是一个vue的项目
  resolve: {                                //路径解析相关的配置
    alias: {                                //alias就是路径的别名
      '@': path.resolve(__dirname, 'src')   //@即表示为绝对路劲 src ；_dirname表示在当前vite...所在的目录
    }
  },
  server: {                                   //服务器的配置
    proxy: {                                  //定义代理规则
      '/api': {                               //若请求路径为api开头
        target: 'http://127.0.0.1:8000',      //则转发路径目标为target所写（Go后端地址）
        changeOrigin: true                    //转发过去后 把host地址改成后端的地址 能够避免后端不处理  
      }
    }
  }
})
