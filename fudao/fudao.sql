
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 数据库： `fudao`
--
CREATE DATABASE IF NOT EXISTS `fudao` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `fudao`;

-- --------------------------------------------------------

--
-- 表的结构 `fudao`
--

CREATE TABLE `fudao` (
  `id` int(11) NOT NULL,
  `ke_id` varchar(64) NOT NULL COMMENT '课程 id',
  `grade` varchar(64) NOT NULL COMMENT '年级',
  `class` varchar(64) NOT NULL COMMENT '学科',
  `title` varchar(255) NOT NULL COMMENT '标题',
  `url` varchar(255) NOT NULL,
  `teacher` varchar(64) NOT NULL COMMENT '教师姓名',
  `price` int(11) NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=MyISAM DEFAULT CHARSET=utf8;


--
-- 转储表的索引
--

--
-- 表的索引 `fudao`
--
ALTER TABLE `fudao`
  ADD PRIMARY KEY (`id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `fudao`
--
ALTER TABLE `fudao`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=95;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
