# gologger

#### 介绍
一个Go中使用的简单的日志处理包，包括Debug，Trace，Info，Warning，Error，Fatal六个级别的日志信息。可以输出日志到控制台、文件、邮件、数据库。文件日志分为wholefile、sizerollingfile和dailyrollingfile：whilefile类型会产生一个整体的日志文件，产生的日志信息不断追加；sizerollingfile类型可以自定义一个日志文件的大小，当到达设定之后会保存该文件，然后创建一个新的文件写入日志,dailyrollingfile会将每天产生的日志保存为一个独立的文件，可以设置要保存的日志的天数。

#### 使用方法
有两种使用方式：xml配置文件方式（推荐）和代码设置方式：

（一）xml配置文件配置

1：在项目中添加xml配置文件（注：由于没有dtd文件的约束，xml文件的格式请严格按照如下及注释描述编写)


    <gologger globalLevel="debug">
        <!-- 配置控制台(此标签及其子标签均为可选项） -->
        <console>
            <!--配置是否输出到控制台，默认为true>
            <enable>false</enable>
            <!--输出到控制台的日志的级别，支持的级别：debug,trace,info,warning,error,fatal>
            <selflevel>debug</selflevel>
            <!--日志格式，ref引用定义的format的id>
            <format ref="first"/>
        </console>
    
        <!-- 配置wholefile，可以在wholefile中配置多个wholefile，将日志输出到不同的wholefile文件。本例配置了两个wholefile。>
        <!-- 此标签可选>
        <wholefiles>
            <!--此标签的path为必选标签>
            <wholefile>
                <!--日志输出文件路径>
                <path>./log/wholeFile.log</path>
                <!--输出到本文件的日志的级别，支持的级别：debug,trace,info,warning,error,fatal>
                <selflevel>info</selflevel>
                <!--日志格式，ref引用定义的format的id>
                <format ref="first"/>
            </wholefile>
            <wholefile>
                <path>./log/wholeFile2.log</path>
                <selflevel>debug</selflevel>
                <format ref="second"/>
            </wholefile>
        </wholefiles>
    
        <!-- 配置dailyrollingfile，可以配置多个，本例配置了两个>
        <!-- 此标签可选>
        <dailyrollingfiles>
            <!--此标签的path为必选标签>
            <dailyrollingfile>
                <!--日志输出文件路径>
                <path>./log/dailyrollingfile/dailyrollingfile.log</path>
                <!--输出到本文件的日志的级别，支持的级别：debug,trace,info,warning,error,fatal>
                <selflevel>trace</selflevel>
                <!--日志格式，ref引用定义的format的id>
                <format ref="first"/>
                <!--保存日志文件的天数，如果为0，则保存过去每天的日志文件>
                <day>2</day>
            </dailyrollingfile>
    
            <dailyrollingfile>
                <path>./log/dailyrollingfile/dailyrollingfile.log</path>
                <selflevel>trace</selflevel>
                <format ref="first"/>
                <day>2</day>
            </dailyrollingfile>
        </dailyrollingfiles>
        
        <!-- 配置sizerollingfile，可以配置多个>
        <!-- 此标签为可选标签>
        <sizerollingfiles>
            <!--此标签的path为必选标签>
            <sizerollingfile>
                <!--日志输出文件路径>
                <path>./log/sizeRollingFile/sizeRollingFile.log</path>
                <!--输出到本文件的日志的级别，支持的级别：debug,trace,info,warning,error,fatal>
                <selflevel>info</selflevel>
                <!--日志格式，ref引用定义的format的id>
                <format ref="first"/>
                <!--一个日志文件的最大容量,单位为Mb>
                <size>1</size>
            </sizerollingfile>
        </sizerollingfiles>
        
        <!-- 配置日志输出到邮件>
        <!-- 此标签为可选标签，如果配置该标签，其子标签host、port、password、from为必选标签>
        <smtplog>
            <!-- 邮件服务器host>
            <host>smtp.qq.com</host>
            <!--监听端口>
            <port>587</port>
            <!--授权码>
            <password>aiaphwysppondfii</password>
            <!--发件邮箱>
            <from>3098086691@qq.com</from>
            <!--日志邮件主题>
            <subject>日志</subject>
            <!--日志邮件发送级别，默认为error及以上>
            <selflevel>debug</selflevel>
            <!--邮件接收者,可以配多个>
            <recipients>
                <!-- 可以配置一个或多个>
                <recipient>3098086691@qq.com</recipient>
                <recipient>2287286256@qq.com</recipient>
            </recipients>
        </smtplog>
    
        <!-- 配置日志输出到数据库>
        <!-- 此标签为可选标签，如果选择，其子标签dirver、username、password、ip、port、databasename为必选标签>
        <databaselog>
            <!--数据库驱动名称>
            <driver>mysql</driver>
            <!--数据库登录用户名>
            <username>root</username>
            <!--数据库密码>
            <password>123456</password>
            <!--数据库服务器ip>
            <ip>127.0.0.1</ip>
            <!--服务器端口>
            <port>3306</port>
            <!--数据库名称>
            <databasename>log_test</databasename>
            <tables>
                <!--日志表，如果没有，会自动创建，可以有多个>
                <table name="newlog">
                    <!-- 输出到该表的日志的级别>
                    <selflevel>debug</selflevel>
                    <!- 该数据表的列，支持的attribute有：date、time、file、func、line、level、msg
                    <attributes>
                        <attribute name="date"/>
                        <attribute name="time"/>
                        <attribute name="msg"/>
                    </attributes>
                </table>
                <table name="anotherlog">
                    <selflevel>info</selflevel>
                    <attributes>
                        <attribute name="file"/>
                        <attribute name="func"/>
                        <attribute name="msg"/>
                    </attributes>
                </table>
            </tables>
        </databaselog>
    
        <formats>
            <!--日志格式，支持的关键词：%Date，%Time，%File，%Func，%Line，%Level，%Msg>
            <format id="first" format="[%Date %Time][%File %Func-%Line] [%Level] %Msg"/>
            <format id="second" format="[%Date %Time][%Level] %Msg"/>
        </formats>
    </gologger>
    
    
  使用方法：
  
  
     在项目中添加xml配置文件，格式如上（一般添加到项目根目录）
     eg:
        logger, err := gologger.GetLoggerByXML("gologger_config.xml")
        if err != nil {
            panic(err)
        }
        defer logger.Flush()
        logger.Debug("debug 测试")
       
    



（二）代码设置方式


1：创建logger对象(必选)

     创建一个日志对象，同时设定全局的日记记录级别（Debug，Trace，Info，Warning，Error，Fatal六种，不区分大小写）
         eg: logger := gologger.NewLogger("fatal")
         
         创建完logger对象后即可使用：logger.Debug("日志信息") 
         输出结果：[DEBUG] 日志信息
   

2：设置输出位置（可选）


     （1）输出日志到控制台
     
         一个刚创建完的logger对象默认输出到控制台。可以通过设置控制是否向控制台输出:
             eg: logger.EnableConsole(true)  开启向控制台输出  
                 logger.EnableConsole(false) 关闭向控制台输出  
          
         根据下面的条目‘3’设置日志输出格式（可选，如不设置，默认输出格式为：[DEBUG] 日志信息）
             
         设置控制台console的私有日志输出级别（可选，如不设置，服从全局日志级别）
             eg：logger.SetConsoleSelfLogLevel(gologger.INFO)


     （2）输出日志到整体文件
     
         创建WholeFile对象（日志直接追加）：
             eg: wholeFile := gologger.NewWholeFile("wholeLog/logtest.log")
             
         根据下面的条目‘3’设置日志输出格式（可选，如不设置，默认输出格式为：[DEBUG] 日志信息）
             
         设置wholeFile的私有日志输出级别（可选，如不设置，服从全局日志级别）：
             eg：wholeFile.SetWholeFileSelfLogLevel(gologger.TRACE)
             
         将wholeFile添加到logger:
             eg: logger.AddWholeFile(wholeFile)
             
         这样即可将日志输出到wholeFile。可以创建多个wholeFile添加到logger。
         

     （3）输出日志到固定大小文件：
     
         创建SizeRollingFile对象（日志追加，如果文件达到设定的maxsize，则进行分割保存，同时创建一个新文件开始写入日志）：
             eg: sizeRollingFile := gologger.NewSizeRollingFile("sizeRollingLog/logtest.log")
             
         根据下面的条目‘3’设置日志输出格式（可选，如不设置，默认输出格式为：[DEBUG] 日志信息）
         
         一个新创建的rollingfile的默认日志大小为10M，可以自己设置大小：
             eg: sizeRollingFile.SetMaxSize(1024) (注：单位是字节）
             
         设置sizeRollingFile的私有日志输出级别（可选，如不设置，服从全局日志级别）：
             eg：sizeRollingFile.SetSizeRollingFileSelfLogLevel(gologger.ERROR)
             
         将sizeRollingFile添加到logger:
             eg: logger.AddSizeRollingFile(sizeRollingFile)
             
                
     （4）输出日志到日记文件
     
         创建DailyRollingFile对象（每天自动保存一个当天的日志文件）
             eg：dailyRollingFile := gologger.NewDailyRollingFile("dailyRollingLog/logtest.log")
             
         根据下面的条目‘3’设置日志输出格式（可选，如不设置，默认输出格式为：[DEBUG] 日志信息）
         
         设定日志保存的天数，如果不设置，默认保存过去7天的日志；如果设置为0，则保存过去所有的日志，不会自动删除；如果输入负数则报错
             eg：dailyRollingFile.SetLogFiletMaxDays(1)
             
         设置dailyRollingFile的私有日志输出级别（可选，如不设置，服从全局日志级别）：
             eg：dailyRollingFile.SetDailyRollingFileSelfLogLevel(gologger.TRACE)
             
         将dailyRollingFile添加到logger：
             eg：logger.AddDailyRollingFile(dailyRollingFile)
             
          
     （5）输出日志到邮件
     
         创建SmtpLog对象，参数依次是：服务器host、授权码passwd、发送邮箱from、开放端口587、邮件主题subject
             eg：smtpLog := gologger.NewSmtpLog("smtp.qq.com", "授权码", "3098086691@qq.com", 587, "日志测试")
         
         设置邮件接收者：参数是数组切片，切片元素是接收邮箱，可以设多个
             eg：smtpLog.SetRecipient([]string{"3098086691@qq.com"})
         
         设置邮件日志发送最低级别（可选），默认是error
             eg：smtpLog.SetSelfMailLevel(gologger.INFO)
         
         将邮件日志对象smtpLog添加到日志记录器中
             eg：logger.AddSmtpLog(smtpLog)
             
            
     （6）输出日志到数据库
        
         创建DatabaseLog对象，参数依次是：数据库驱动名，数据库系统账号，密码，数据库服务器ip，端口号，数据库名称
             eg：databaseLog := gologger.NewDatabaseLog("mysql", "root", "zhangqi", "127.0.0.1", 3306, "log_test")
        
         创建日志表对象，第一个参数是“表名”，后续参数是日至标的数据列
             eg：logTable := gologger.NewLogTable("log", gologger.Date,gologger.Time, gologger.File, gologger.Func, gologger.Line, gologger.Level, gologger.Msg)
         
         设置该日志表记录日志的私有级别
             eg：logTable.SetTableSelfLevel(gologger.INFO)（可选，如不设置，服从全局日志级别）
        
         将日志表对象添加到数据库，数据库对象可以添加多个日志表对象
             eg：databaseLog.AddTable(logTable)
        
         将数据库对象添加到日志记录器logger
             eg：logger.AddDatabaseLog(databaseLog)
        


3：设置日志格式format（可选）


    每一个wholeFile，sizeRollingFile、dailyRollingFile和console都可以设置日志输出的格式:
    
       创建format对象：
          eg：format := gologger.DefaultFormat() （注：获取默认日志格式对象：[2020-07-27 14:09:52][Fatal][D:/Go_WorkSpace/src/gologger/main/main.go：main.test：34]fatal测试：李四）
              format := gologger.NewFormat("[%Date %Time][%File %Func-%Line] [%Level] %Msg")
          
          注：支持的日志格式关键字
          keyWord |  show
          --------------------
          %Date   |  2006-01-02
          %Time   |  15:04:05
          %Level  |  Fatal
          %File   |  D:/Go_WorkSpace/src/gologger/main/main.go
          %Func   |  main.test
          %Line   |  33
          %Msg    |  debug测试
          

       将创建的format赋值给wholeFile
          eg: wholeFile.SetWholeFileFormat(format)
          
       将创建的format赋值给sizeRollingFile
          eg: sizeRollingFile.SetSizeRollingFileFormat(format)
          
       将创建的format赋值给dailyRollingFile
          eg: dailyRollingFile.SetDailyRollingFileFormat(format)
          
       将创建的format赋值给控制台console
          eg: logger.SetConsoleFormat(format)

   
5: 关闭日志(必选)

    eg：defer logger.Flush()
       
#### 版本
    v2.0.1
        *改变邮件日志格式：从html的table——>ul/li
        *扩大Logger的访问权限：logger(包可见)——>Logger(全局可见)
    v2.0.0
        *增加使用xml配置文件的配置方式，简化配置过程
        *修复若干bug

    v1.5.0
        *增加数据库日志功能
            
    v1.4.0
        *增加邮件日志功能
            
    v1.3.1
        *增加设定保存日志的天数功能
            
    v1.3.0
        *增加每天日志自动分割功能
        *修复日志异步写入bug

    v1.2.0
        *增加自定义日志格式功能，支持日志记录格式的私人定制
        
    v1.1.1
        *修复了一个bug
    
    v1.1.0
        *完成日志的异步输出
        *可以向控制台console、整体日志file和rollingfile日志分别独立输出
        *可以给控制台console、整体日志file和rollingfile日志设置不同的日志记录格式
        *日志记录格式暂时不支持私人定制，只有内置的默认格式DefaultFormat()
        

#### 联系
     邮箱:3098086691@qq.com
     Q Q: 3098086691